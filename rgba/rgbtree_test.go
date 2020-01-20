package rgba

import (
	"fmt"
	"image/color"
	"image/color/palette"
	"math/rand"
	"testing"

	"github.com/shabbyrobe/imgx/testimg"
)

func TestRGBNN(t *testing.T) {
	const iter = 10000

	type paletteCase struct {
		name string
		pal  color.Palette
	}

	palettes := []paletteCase{
		{"websafe32", color.Palette(palette.WebSafe[:32])},
		{"websafe216", color.Palette(palette.WebSafe)},
	}

	for i := 1; i <= 256; i++ {
		rng := rand.New(rand.NewSource(0))
		pal := make(color.Palette, i)
		for c := 0; c < i; c++ {
			pal[c] = testimg.RandRGBA(rng)
		}
		palettes = append(palettes, paletteCase{fmt.Sprintf("rand%d", i), pal})
	}

	var same, diff int
	for _, pc := range palettes {
		rpal := ConvertPalette(pc.pal)

		node := rgbTreeBuild(rpal)
		rng := rand.New(rand.NewSource(0))

		for i := 0; i < iter; i++ {
			col := testimg.RandRGBA(rng)
			idx1 := rgbNearestEuclideanIndex(rpal, col)

			cnv1 := rpal[idx1]
			cnv2 := node.NearestRGBAColor(col)
			if cnv1 != cnv2 {
				diff++
				dist1 := (0 +
					sqDiff8(cnv1.R, col.R) +
					sqDiff8(cnv1.G, col.G) +
					sqDiff8(cnv1.B, col.B))

				dist2 := (0 +
					sqDiff8(cnv2.R, col.R) +
					sqDiff8(cnv2.G, col.G) +
					sqDiff8(cnv2.B, col.B))

				if dist1 != dist2 {
					t.Fatal(pc.name, i, "col:", col, "expected:", cnv1, "found:", cnv2,
						"eucliddist:", dist1, "nndist:", dist2)
				}

			} else {
				same++

				// It only makes sense to test the index if the result is the same
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

var BenchRGBNodeResult *rgbNode

func BenchmarkRGBTreeBuild(b *testing.B) {
	pal := ConvertPalette(palette.WebSafe)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BenchRGBNodeResult = rgbTreeBuild(pal)
	}
}
