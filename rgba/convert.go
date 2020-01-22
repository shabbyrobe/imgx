package rgba

import (
	"image"
	"image/color"
)

type rgbaAtImage interface {
	image.Image
	RGBAAt(x, y int) color.RGBA
}

var _ rgbaAtImage = &image.RGBA{}

// CastFromBytes casts a raw byte array into a slice of color.RGBA values. After calling this,
// data is no longer safe to use independently.
//
// The length of data must be a multiple of 4.
func CastFromBytes(data []byte) ([]color.RGBA, error) {
	return castFromBytes(data)
}

// CastToBytes casts a slice of color.RGBA values into a byte slice, avoiding copies if
// possible. After calling this, colors is no longer safe to use independently.
func CastToBytes(colors []color.RGBA) ([]byte, error) {
	return castToBytes(colors)
}

// Convert an image.Image into an *image.RGBA.
//
// This will attempt to cast the image first, and if that succeeds, 'copied' will be
// false. If you require a copy of the image, call CloneDeep() on the output:
//
//	var img image.Image
//	rimg, copied := rgba.Convert(img)
//	if !copied {
//		rimg = rimg.CloneDeep()
//	}
//
// If 'copied' is false and you do not clone the output, the original image is no longer
// safe to use.
//
// Most of the image formats from the stdlib have a "fast" conversion path in here, but if
// one does not exist it should be added. If a fast path is unavailable, there are slow
// paths that attempt to use image.RGBAAt(), then finally image.At().
//
func Convert(img image.Image) (out *Image, copied bool) {
	switch img := img.(type) {
	case *Image:
		return img, false
	case *image.CMYK:
		return convertCMYKToRGBA(img), true
	case *image.NRGBA:
		return convertNRGBAToRGBA(img), true
	case *image.NRGBA64:
		return convertNRGBA64ToRGBA(img), true
	case *image.Paletted:
		return convertPalettedToRGBA(img), true
	case *image.RGBA:
		return convertRGBAToRGBA(img)
	case *image.RGBA64:
		return convertRGBA64ToRGBA(img), true
	case *image.YCbCr:
		return convertYCbCrToRGBA(img), true
	case rgbaAtImage:
		return convertRGBAAtToRGBA(img), true
	default:
		return convertImageToRGBA(img), true
	}
}

func convertCMYKToRGBA(img *image.CMYK) *Image {
	size := img.Bounds().Size()
	out := New(size)
	inPix, outVals := img.Pix, out.Vals

	for i, j := 0, 0; i < len(inPix); i, j = i+4, j+1 {
		// Seems like it might be unnecessary to go from 8-bit to 16-bit to
		// 8-bit again, but I'm not quite sure yet and haven't looked further:
		w := 0xffff - uint32(inPix[i+3])*0x101
		outVals[j].R = uint8(((0xffff - uint32(inPix[i+0])*0x101) * w / 0xffff) >> 8)
		outVals[j].G = uint8(((0xffff - uint32(inPix[i+1])*0x101) * w / 0xffff) >> 8)
		outVals[j].B = uint8(((0xffff - uint32(inPix[i+2])*0x101) * w / 0xffff) >> 8)
		outVals[j].A = 0xff
	}

	return out
}

func convertPalettedToRGBA(img *image.Paletted) *Image {
	// FIXME: if img.Palette is too big, it might be worth just using
	// the generic converter:
	pal := make([]color.RGBA, len(img.Palette))

	for idx, col := range img.Palette {
		r, g, b, a := col.RGBA()
		pal[idx] = color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
	}

	size := img.Bounds().Size()
	out := New(size)
	inPix, outVals := img.Pix, out.Vals
	for i, c := range inPix {
		outVals[i] = pal[c]
	}
	return out
}

func convertYCbCrToRGBA(img *image.YCbCr) *Image {
	bounds := img.Bounds()
	size := bounds.Size()
	vals := make([]color.RGBA, size.X*size.Y)

	var pix int
	var accumYOffset int
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			// image:YCbCr.YOffset():
			var yOffset = accumYOffset + x
			var cOffset int

			// {{{ image.YCbCr.COffset():
			switch img.SubsampleRatio {
			case image.YCbCrSubsampleRatio422:
				cOffset = (y)*img.CStride + (x / 2)
			case image.YCbCrSubsampleRatio420:
				cOffset = (y/2)*img.CStride + (x / 2)
			case image.YCbCrSubsampleRatio440:
				cOffset = (y/2)*img.CStride + (x)
			case image.YCbCrSubsampleRatio411:
				cOffset = (y)*img.CStride + (x / 4)
			case image.YCbCrSubsampleRatio410:
				cOffset = (y/2)*img.CStride + (x / 4)
			default:
				cOffset = (y)*img.CStride + (x)
			}
			// }}}

			y, cb, cr := img.Y[yOffset], img.Cb[cOffset], img.Cr[cOffset]

			// {{{ image.YCbCrToRGB():
			yy1 := int32(y) * 0x10101
			cb1 := int32(cb) - 128
			cr1 := int32(cr) - 128

			r := yy1 + 91881*cr1
			if uint32(r)&0xff000000 == 0 {
				r >>= 16
			} else {
				r = ^(r >> 31)
			}

			g := yy1 - 22554*cb1 - 46802*cr1
			if uint32(g)&0xff000000 == 0 {
				g >>= 16
			} else {
				g = ^(g >> 31)
			}

			b := yy1 + 116130*cb1
			if uint32(b)&0xff000000 == 0 {
				b >>= 16
			} else {
				b = ^(b >> 31)
			}
			// }}}

			vals[pix] = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xff}
			pix++
		}
		accumYOffset += img.YStride
	}

	return &Image{Size: size, Stride: size.X, Vals: vals}
}

func convertNRGBA64ToRGBA(img *image.NRGBA64) *Image {
	size := img.Bounds().Size()
	out := New(size)
	inPix, outVals := img.Pix, out.Vals
	inLen := len(inPix)

	for ip, op := 0, 0; ip < inLen; ip, op = ip+8, op+1 {
		a := (uint32(inPix[ip+6]) << 8) | uint32(inPix[ip+7])

		outVals[op] = color.RGBA{
			R: uint8(((uint32(inPix[ip+0]) << 8) | uint32(inPix[ip+1])*a/0xffff) >> 8),
			G: uint8(((uint32(inPix[ip+2]) << 8) | uint32(inPix[ip+3])*a/0xffff) >> 8),
			B: uint8(((uint32(inPix[ip+4]) << 8) | uint32(inPix[ip+5])*a/0xffff) >> 8),
			A: uint8(a >> 8),
		}
	}

	return out
}

func convertNRGBAToRGBA(img *image.NRGBA) *Image {
	size := img.Bounds().Size()
	out := New(size)
	inPix, outVals := img.Pix, out.Vals

	for ip, op := 0, 0; ip < len(inPix); ip, op = ip+4, op+1 {
		a := uint32(inPix[ip+3])
		outVals[op] = color.RGBA{
			R: uint8(uint32(inPix[ip+0]) * a / 0xff),
			G: uint8(uint32(inPix[ip+1]) * a / 0xff),
			B: uint8(uint32(inPix[ip+2]) * a / 0xff),
			A: uint8(a),
		}
	}

	return out
}

func convertRGBA64ToRGBA(img *image.RGBA64) *Image {
	size := img.Bounds().Size()
	out := New(size)
	inPix, outVals := img.Pix, out.Vals

	for ip, op := 0, 0; ip < len(inPix); ip, op = ip+8, op+1 {
		// RGBA64 stores pixels in big-endian pairs. We only need the big end:
		outVals[op] = color.RGBA{
			R: inPix[ip+0],
			G: inPix[ip+2],
			B: inPix[ip+4],
			A: inPix[ip+6],
		}
	}

	return out
}

// convertRGBtToRGBASlow is used if we can't use the faster type-pun version
// found in convert_fast.go.
func convertRGBAToRGBASlow(img *image.RGBA) (out *Image, copied bool) {
	size := img.Bounds().Size()

	out = New(size)
	inPix, outVals := img.Pix, out.Vals
	inLen := len(inPix)

	for ip, op := 0, 0; ip < inLen; ip, op = ip+4, op+1 {
		outVals[op] = color.RGBA{
			R: inPix[ip+0],
			G: inPix[ip+1],
			B: inPix[ip+2],
			A: inPix[ip+3],
		}
	}

	return out, true
}

// convertRGBAAtToRGBA is hopefully a less grim fallback slow-path than the
// CPU-warmer convertImageToRGBA.
func convertRGBAAtToRGBA(img rgbaAtImage) *Image {
	bounds := img.Bounds()
	size := bounds.Size()
	vals := make([]color.RGBA, size.X*size.Y)

	var pix int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			vals[pix] = img.RGBAAt(x, y)
			pix++
		}
	}

	return &Image{Size: size, Stride: size.X, Vals: vals}
}

func convertImageToRGBA(img image.Image) *Image {
	bounds := img.Bounds()
	size := bounds.Size()
	vals := make([]color.RGBA, size.X*size.Y)

	var pix int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			vals[pix] = color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
			pix++
		}
	}

	return &Image{Size: size, Stride: size.X, Vals: vals}
}
