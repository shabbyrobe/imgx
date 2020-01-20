package rgba

import (
	"image/color"
)

var DefaultIndexer = NewRGBATreeIndexer()

// Palette of color.RGBA values.
type Palette []color.RGBA

// ColorPalette returns the rgba.Palette as a color.Palette containing only
// color.RGBA values.
func (p Palette) ColorPalette() color.Palette {
	out := make(color.Palette, len(p))
	for idx, c := range p {
		out[idx] = c
	}
	return out
}

// Index the palette using the default indexer:
func (p Palette) Index() Index {
	return DefaultIndexer.IndexRGBAPalette(p)
}

// ConvertPalette from a color.Palette containing arbitrary color types to an
// rgba.Palette.
func ConvertPalette(pal color.Palette) (out Palette) {
	out = make(Palette, len(pal))

	for idx, c := range pal {
		switch c := c.(type) {
		case color.RGBA:
			out[idx] = c

		case color.RGBA64:
			out[idx] = color.RGBA{
				R: uint8(c.R >> 8),
				G: uint8(c.G >> 8),
				B: uint8(c.B >> 8),
				A: uint8(c.A >> 8),
			}

		case color.NRGBA:
			a32 := uint32(c.A)
			out[idx] = color.RGBA{
				R: uint8(uint32(c.R) * a32 / 0xff),
				G: uint8(uint32(c.G) * a32 / 0xff),
				B: uint8(uint32(c.B) * a32 / 0xff),
				A: c.A,
			}

		default:
			cr, cg, cb, ca := c.RGBA()
			out[idx] = color.RGBA{
				R: uint8(cr >> 8),
				G: uint8(cg >> 8),
				B: uint8(cb >> 8),
				A: uint8(ca >> 8),
			}
		}
	}

	return out
}

// NormalizePalette from a color.Palette containing arbitrary color types to a
// color.Palette containing only color.RGBA values.
func NormalizePalette(pal color.Palette) (out color.Palette) {
	out = make(color.Palette, len(pal))

	for idx, c := range pal {
		switch c := c.(type) {
		case color.RGBA:
			out[idx] = c

		case color.RGBA64:
			out[idx] = color.RGBA{
				R: uint8(c.R >> 8),
				G: uint8(c.G >> 8),
				B: uint8(c.B >> 8),
				A: uint8(c.A >> 8),
			}

		case color.NRGBA:
			a32 := uint32(c.A)
			out[idx] = color.RGBA{
				R: uint8(uint32(c.R) * a32 / 0xff),
				G: uint8(uint32(c.G) * a32 / 0xff),
				B: uint8(uint32(c.B) * a32 / 0xff),
				A: c.A,
			}

		default:
			cr, cg, cb, ca := c.RGBA()
			out[idx] = color.RGBA{
				R: uint8(cr >> 8),
				G: uint8(cg >> 8),
				B: uint8(cb >> 8),
				A: uint8(ca >> 8),
			}
		}
	}
	return out
}
