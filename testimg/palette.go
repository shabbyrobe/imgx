package testimg

import (
	"image/color"
	"math/rand"
)

func RandPalette(rng *rand.Rand, sz int) color.Palette {
	if sz <= 0 || sz > 256 {
		panic("size must be between 1 and 256")
	}
	if rng == nil {
		rng = defaultRNG
	}

	pal := make(color.Palette, sz)
	for i := 0; i < sz; i++ {
		next := uint32(rng.Uint32())
		r, g, b := uint8(next<<16), uint8(next<<8), uint8(next)
		pal[i] = color.RGBA{R: r, G: g, B: b, A: 0xff}
	}
	return pal
}

func RandRGBA(rng *rand.Rand) color.RGBA {
	v := rng.Uint32()
	col := color.RGBA{uint8(v >> 24), uint8(v >> 16), uint8(v >> 8), uint8(v)}

	if col.A < col.R {
		col.A = col.R
	}
	if col.A < col.G {
		col.A = col.G
	}
	if col.A < col.B {
		col.A = col.B
	}
	return col
}

func RandRGBA64(rng *rand.Rand) color.RGBA64 {
	v := rng.Uint64()
	col := color.RGBA64{uint16(v >> 48), uint16(v >> 32), uint16(v >> 16), uint16(v)}

	if col.A < col.R {
		col.A = col.R
	}
	if col.A < col.G {
		col.A = col.G
	}
	if col.A < col.B {
		col.A = col.B
	}
	return col
}
