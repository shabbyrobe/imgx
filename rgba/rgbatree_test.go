package rgba

import (
	"fmt"
	"image/color"
	"image/color/palette"
	"math/rand"
	"strings"
	"testing"

	"github.com/shabbyrobe/imgx/testimg"
)

func TestRGBATreeIndexer(t *testing.T) {
	const iter = 10000

	var rng = rand.New(rand.NewSource(0))

	var same, diff int

	for _, pc := range paletteCases(rng) {
		pal := pc.pal
		rpal := ConvertPalette(pc.pal)

		node := rgbaTreeBuild(rpal)

		for i := 0; i < iter; i++ {
			col := testimg.RandRGBA(rng)
			cnv1 := pal.Convert(col).(color.RGBA)
			cnv2 := node.NearestRGBAColor(col)
			if cnv1 != cnv2 {
				diff++
				dist1 := (0 +
					sqDiff8(cnv1.R, col.R) +
					sqDiff8(cnv1.G, col.G) +
					sqDiff8(cnv1.B, col.B) +
					sqDiff8(cnv1.A, col.A))

				dist2 := (0 +
					sqDiff8(cnv2.R, col.R) +
					sqDiff8(cnv2.G, col.G) +
					sqDiff8(cnv2.B, col.B) +
					sqDiff8(cnv2.A, col.A))

				if dist1 != dist2 {
					t.Fatal(pc.name, i, "col:", col, "expected:", cnv1, "found:", cnv2,
						"eucliddist:", dist1, "nndist:", dist2)
				}

			} else {
				same++

				// It only makes sense to test the index if the result is the same
				idx1 := pal.Index(col)
				idx2 := node.NearestRGBAIndex(col)
				if idx1 != idx2 {
					t.Fatal(pc.name, i, "idx:", idx1, "!=", idx2)
				}
			}
		}
	}

	if same == 0 || diff == 0 {
		t.Fatal("suspicious results", "same:", same, "diff:", diff)
	}
}

var BenchNodeResult *rgbaNode

func BenchmarkRGBATreeBuild(b *testing.B) {
	pal := ConvertPalette(palette.WebSafe)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchNodeResult = rgbaTreeBuild(pal)
	}
}

func dumpNode(kd *rgbaNode, kind string, depth int) string {
	out := ""
	ind := strings.Repeat("  ", depth)
	if kind != "" {
		kind += " "
	}
	out += ind + kind + fmt.Sprintf("%s %v", kd.axis, kd.RGBA()) + "\n"
	if kd.left != nil {
		out += dumpNode(kd.left, "L", depth+1)
	}
	if kd.right != nil {
		out += dumpNode(kd.right, "R", depth+1)
	}
	return out
}
