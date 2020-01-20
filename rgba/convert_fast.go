//+build !purego

package rgba

import (
	"fmt"
	"image"
	"image/color"
	"reflect"
	"unsafe"
)

func castBytes(data []byte) ([]color.RGBA, error) {
	if len(data)%4 != 0 {
		return nil, fmt.Errorf("rgba: raw RGBA data must be a multiple of 4")
	}

	// TYPE PUNNING. YOU FIEND. THIS IS TOTAL EVIL. UNSPEAKABLE EVIL. POSSIBLY TOO EVIL.
	// EXTREME DANGER. If explosions occur, I am sorry. We will have to go back to slow
	// mode if they occur.
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&data))
	header.Len /= 4
	header.Cap /= 4
	return *(*[]color.RGBA)(unsafe.Pointer(&header)), nil
}

var colorVal color.RGBA

func convertRGBAToRGBA(img *image.RGBA) (out *Image, copied bool) {
	if unsafe.Sizeof(colorVal) != 4 {
		return convertRGBAToRGBASlow(img)
	}

	vals, err := CastBytes(img.Pix)
	if err != nil {
		panic(err)
	}

	size := img.Bounds().Size()
	return &Image{Size: size, Stride: size.X, Vals: vals}, false
}
