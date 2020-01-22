//+build !purego

package rgba

import (
	"fmt"
	"image"
	"image/color"
	"reflect"
	"unsafe"
)

func castFromBytes(data []byte) ([]color.RGBA, error) {
	if len(data)%4 != 0 {
		return nil, fmt.Errorf("rgba: raw RGBA data must be a multiple of 4")
	}
	// EVIL:
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&data))
	header.Len /= 4
	header.Cap /= 4
	return *(*[]color.RGBA)(unsafe.Pointer(&header)), nil
}

func castToBytes(colors []color.RGBA) (data []byte, err error) {
	// EVIL:
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&colors))
	header.Len *= 4
	header.Cap *= 4
	return *(*[]byte)(unsafe.Pointer(&header)), nil
}

var colorVal color.RGBA

func convertRGBAToRGBA(img *image.RGBA) (out *Image, copied bool) {
	if unsafe.Sizeof(colorVal) != 4 {
		return convertRGBAToRGBASlow(img)
	}

	vals, err := CastFromBytes(img.Pix)
	if err != nil {
		panic(err)
	}

	size := img.Bounds().Size()
	return &Image{Size: size, Stride: size.X, Vals: vals}, false
}
