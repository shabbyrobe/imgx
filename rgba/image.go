package rgba

import (
	"image"
	"image/color"
)

// Image containing pixels of color.RGBA values.
//
// This exists mostly as a readability measure over image.RGBA, though it seems to
// benchmark a bit faster in whole-system tests. image.RGBA.Pix's int-based offsets tend
// to introduce a lot of extra noise to some of the image algorithms built around it
// (which I already find hard enough to read).
//
type Image struct {
	Size   image.Point
	Stride int
	Vals   []color.RGBA
}

var _ image.Image = &Image{}

func New(size image.Point) *Image {
	return &Image{
		Size:   size,
		Stride: size.X,
		Vals:   make([]color.RGBA, size.X*size.Y),
	}
}

func NewWithVals(size image.Point, vals []color.RGBA) *Image {
	if len(vals) != size.X*size.Y {
		panic("rgba: size did not match len(vals)")
	}
	return &Image{Size: size, Stride: size.X, Vals: vals}
}

func (p *Image) CloneDeep() *Image {
	vals := make([]color.RGBA, len(p.Vals))
	copy(vals, p.Vals)
	return &Image{Size: p.Size, Stride: p.Stride, Vals: vals}
}

func (p *Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (p *Image) Bounds() image.Rectangle {
	return image.Rectangle{Max: p.Size}
}

func (p *Image) At(x, y int) color.Color {
	return p.RGBAAt(x, y)
}

func (p *Image) PixOffset(x, y int) int {
	return y*p.Stride + x
}

func (p *Image) RGBAAt(x, y int) (c color.RGBA) {
	if x >= p.Size.X || y >= p.Size.Y {
		return c
	}
	return p.Vals[y*p.Stride+x]
}

func (p *Image) Set(x, y int, c color.Color) {
	if x >= p.Size.X || y >= p.Size.Y {
		return
	}

	switch c := c.(type) {
	case color.RGBA:
		p.Vals[y*p.Stride+x] = c

	case color.RGBA64:
		p.Vals[y*p.Stride+x] = color.RGBA{
			R: uint8(c.R >> 8),
			G: uint8(c.G >> 8),
			B: uint8(c.B >> 8),
			A: uint8(c.A >> 8),
		}

	case color.NRGBA:
		a32 := uint32(c.A)
		p.Vals[y*p.Stride+x] = color.RGBA{
			R: uint8(uint32(c.R) * a32 / 0xff),
			G: uint8(uint32(c.G) * a32 / 0xff),
			B: uint8(uint32(c.B) * a32 / 0xff),
			A: c.A,
		}

	default:
		r, g, b, a := c.RGBA()
		p.Vals[y*p.Stride+x] = color.RGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: uint8(a >> 8),
		}
	}
}

func (p *Image) SetRGBA(x, y int, c color.RGBA) {
	if x >= p.Size.X || y >= p.Size.Y {
		return
	}
	p.Vals[y*p.Stride+x] = c
}
