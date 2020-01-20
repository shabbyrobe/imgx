package testimg

import (
	"image"
	"image/color"
	"math/rand"
)

type Recipe interface {
	RGBA(*rand.Rand) *image.RGBA
	RGBA64(*rand.Rand) *image.RGBA64
	NRGBA(*rand.Rand) *image.NRGBA
	NRGBA64(*rand.Rand) *image.NRGBA64
	Paletted(*rand.Rand, color.Palette) *image.Paletted
	YCbCr(*rand.Rand) *image.YCbCr
	CMYK(*rand.Rand) *image.CMYK
}
