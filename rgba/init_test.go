package rgba

import (
	"fmt"
	"image/color"
	"image/color/palette"
	"math/rand"

	"github.com/shabbyrobe/imgx/testimg"
)

type paletteCase struct {
	name string
	pal  color.Palette
}

func paletteCases(rng *rand.Rand) []paletteCase {
	palettes := []paletteCase{
		{"websafe32", color.Palette(palette.WebSafe[:32])},
		{"websafe216", color.Palette(palette.WebSafe)},
	}

	if rng == nil {
		rng = rand.New(rand.NewSource(0))
	}

	for i := 1; i <= 256; i++ {
		pal := make(color.Palette, i)
		for c := 0; c < i; c++ {
			pal[c] = testimg.RandRGBA(rng)
		}
		palettes = append(palettes, paletteCase{fmt.Sprintf("rand%d", i), pal})
	}

	return palettes
}

func rgbNearestEuclideanIndex(p Palette, c color.RGBA) int {
	cr, cg, cb := c.R, c.G, c.B
	ret, bestSum := 0, uint32(1<<32-1)

	for i, v := range p {
		vr, vg, vb, _ := v.RGBA()
		sum := 0 +
			sqDiff8(cr, uint8(vr>>8)) +
			sqDiff8(cg, uint8(vg>>8)) +
			sqDiff8(cb, uint8(vb>>8))

		if sum < bestSum {
			if sum == 0 {
				return i
			}
			ret, bestSum = i, sum
		}
	}

	return ret
}

func rgbNearestEuclidean(p Palette, c color.RGBA) color.RGBA {
	return p[rgbNearestEuclideanIndex(p, c)]
}
