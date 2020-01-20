package rgba

import (
	"fmt"
	"image/color"
	"math/rand"
	"testing"

	"github.com/shabbyrobe/imgx/testimg"
)

var max int64

func TestRGBPrecacheIndexer(t *testing.T) {
	const iter = 10000

	// 3 bits of error for 3 color channels, squared:
	const maxError = (1 << (3 + 3 + 3)) * (1 << (3 + 3 + 3))

	var same, diff int
	var indexer = NewRGBPrecacheIndexer(nil)

	defer func() {
		fmt.Println(max)
	}()

	rng := rand.New(rand.NewSource(0))

	for _, pc := range paletteCases(rng) {
		pal := pc.pal
		rpal := ConvertPalette(pc.pal)
		node := indexer.IndexRGBAPalette(rpal)

		for i := 0; i < iter; i++ {
			col := testimg.RandRGBA(rng)
			cnv1 := pal.Convert(col).(color.RGBA)
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

				diff := int64(dist1) - int64(dist2)
				if diff < 0 {
					diff = -diff
				}

				if diff > max {
					max = diff
				}

				if diff > maxError {
					t.Fatal(pc.name, i, "col:", col, "expected:", cnv1, "found:", cnv2,
						"eucliddist:", dist1, "nndist:", dist2, "diff:", diff, "maxerr:", maxError)
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

var (
	BenchIndexResult       Index
	BenchSearchResult      int
	BenchSearchColorResult color.RGBA
)

func BenchmarkRGBPrecacheIndexRGBAPalette(b *testing.B) {
	rng := rand.New(rand.NewSource(0))

	pal := testimg.RandPalette(rng, 256)
	rpal := ConvertPalette(pal)
	idxr := NewRGBPrecacheIndexer(nil)

	for i := 0; i < b.N; i++ {
		BenchIndexResult = idxr.IndexRGBAPalette(rpal)
	}
}

func BenchmarkRGBPrecacheRGBASearch(b *testing.B) {
	rng := rand.New(rand.NewSource(0))

	pal := testimg.RandPalette(rng, 256)
	rpal := ConvertPalette(pal)
	idx := NewRGBPrecacheIndexer(nil).IndexRGBAPalette(rpal)

	colCnt := 10000
	cols := make([]color.RGBA, colCnt)
	for i := 0; i < colCnt; i++ {
		cols[i] = testimg.RandRGBA(rng)
	}

	b.Run("", func(b *testing.B) {
		for i, j := 0, 0; i < b.N; i++ {
			BenchSearchResult = idx.NearestRGBAIndex(cols[j])
			j++
			if j >= colCnt {
				j = 0
			}
		}
	})

	b.Run("", func(b *testing.B) {
		for i, j := 0, 0; i < b.N; i++ {
			BenchSearchColorResult = idx.NearestRGBAColor(cols[j])
			j++
			if j >= colCnt {
				j = 0
			}
		}
	})
}
