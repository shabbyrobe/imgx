package testimg

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
)

type RandBlocks struct {
	W, H           int
	BlockW, BlockH int
}

func (r *RandBlocks) ensureValid() {
	if r.BlockW <= 0 {
		r.BlockW = 1
	}
	if r.BlockH <= 0 {
		r.BlockH = 1
	}
	if r.W <= 0 || r.H <= 0 {
		panic("testimg: missing size")
	}
}

func (r RandBlocks) gen(rng *rand.Rand, makeSet func() func(x, y int)) {
	r.ensureValid()
	if rng == nil {
		rng = defaultRNG
	}

	for x := 0; x < r.W; x += r.BlockW {
		for y := 0; y < r.H; y += r.BlockH {
			set := makeSet()
			for bx := 0; bx < r.BlockW; bx++ {
				for by := 0; by < r.BlockH; by++ {
					if x+bx < r.W && y+by < r.H {
						set(x+bx, y+by)
					}
				}
			}
		}
	}
}

func (r RandBlocks) Paletted(rng *rand.Rand, palette color.Palette) *image.Paletted {
	if palette == nil {
		palette = defaultRandPalette
	}
	if rng == nil {
		rng = defaultRNG
	}
	var img = image.NewPaletted(image.Rect(0, 0, r.W, r.H), palette)
	var makeSet = func() func(x, y int) {
		col := uint8(rng.Intn(len(palette)))
		return func(x, y int) {
			img.Pix[r.W*y+x] = col
		}
	}
	r.gen(rng, makeSet)
	return img
}

func (r RandBlocks) RGBA(rng *rand.Rand) *image.RGBA {
	if rng == nil {
		rng = defaultRNG
	}
	var img = image.NewRGBA(image.Rect(0, 0, r.W, r.H))
	var makeSet = func() func(x, y int) {
		v := rng.Uint32()
		cr, cg, cb := uint8(v>>16), uint8(v>>8), uint8(v)
		return func(x, y int) {
			i := (r.W*y + x) * 4
			img.Pix[i+0] = cr
			img.Pix[i+1] = cg
			img.Pix[i+2] = cb
			img.Pix[i+3] = 0xFF
		}
	}
	r.gen(rng, makeSet)
	return img
}

func (r RandBlocks) RGBA64(rng *rand.Rand) *image.RGBA64 {
	if rng == nil {
		rng = defaultRNG
	}
	var img = image.NewRGBA64(image.Rect(0, 0, r.W, r.H))
	var makeSet = func() func(x, y int) {
		v := rng.Uint32()
		return func(x, y int) {
			i := (r.W*y + x) * 8
			img.Pix[i+0] = uint8(v >> 16)
			img.Pix[i+1] = uint8(v >> 16)
			img.Pix[i+2] = uint8(v >> 8)
			img.Pix[i+3] = uint8(v >> 8)
			img.Pix[i+4] = uint8(v)
			img.Pix[i+5] = uint8(v)
			img.Pix[i+6] = 0xff
			img.Pix[i+7] = 0xff
		}
	}
	r.gen(rng, makeSet)
	return img
}

func (r RandBlocks) NRGBA(rng *rand.Rand) *image.NRGBA {
	if rng == nil {
		rng = defaultRNG
	}
	var img = image.NewNRGBA(image.Rect(0, 0, r.W, r.H))
	var makeSet = func() func(x, y int) {
		v := rng.Uint32()
		cr, cg, cb := uint8(v>>16), uint8(v>>8), uint8(v)
		return func(x, y int) {
			i := (r.W*y + x) * 4
			img.Pix[i+0] = cr
			img.Pix[i+1] = cg
			img.Pix[i+2] = cb
			img.Pix[i+3] = 0xFF
		}
	}
	r.gen(rng, makeSet)
	return img
}

func (r RandBlocks) NRGBA64(rng *rand.Rand) *image.NRGBA64 {
	if rng == nil {
		rng = defaultRNG
	}
	var img = image.NewNRGBA64(image.Rect(0, 0, r.W, r.H))
	var makeSet = func() func(x, y int) {
		v := rng.Uint32()
		return func(x, y int) {
			i := (r.W*y + x) * 8
			img.Pix[i+0] = uint8(v >> 16)
			img.Pix[i+1] = uint8(v >> 16)
			img.Pix[i+2] = uint8(v >> 8)
			img.Pix[i+3] = uint8(v >> 8)
			img.Pix[i+4] = uint8(v)
			img.Pix[i+5] = uint8(v)
			img.Pix[i+6] = 0xff
			img.Pix[i+7] = 0xff
		}
	}
	r.gen(rng, makeSet)
	return img
}

func (r RandBlocks) CMYK(rng *rand.Rand) *image.CMYK {
	if rng == nil {
		rng = defaultRNG
	}
	var img = image.NewCMYK(image.Rect(0, 0, r.W, r.H))
	var makeSet = func() func(x, y int) {
		v := rng.Uint32()
		cc, cm, cy, ck := uint8(v>>24), uint8(v>>16), uint8(v>>8), uint8(v)
		return func(x, y int) {
			i := (r.W*y + x) * 4
			img.Pix[i+0] = cc
			img.Pix[i+1] = cm
			img.Pix[i+2] = cy
			img.Pix[i+3] = ck
		}
	}
	r.gen(rng, makeSet)
	return img
}

func (r RandBlocks) YCbCr(rng *rand.Rand) *image.YCbCr {
	// No way to create a YCbCr image without lots of spelunking and reading maths. It's
	// late, and this works. This only produces a 4:2:0 image; ultimately I'll need to
	// work out how to encode these things from RGB to hit all those bases but not
	// tonight.
	img := r.RGBA(rng)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100}); err != nil {
		panic(err)
	}
	enc, err := jpeg.Decode(&buf)
	if err != nil {
		panic(err)
	}
	ycb, ok := enc.(*image.YCbCr)
	if !ok {
		panic("could not produce YCbCr")
	}
	return ycb
}
