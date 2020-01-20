//+build purego

package rgba

import (
	"fmt"
	"image"
	"image/color"
)

func castBytes(data []byte) ([]color.RGBA, error) {
	dlen := len(data)
	if dlen%4 != 0 {
		return nil, fmt.Errorf("rgba: raw RGBA data must be a multiple of 4")
	}

	out := make([]color.RGBA, dlen)
	for ip, op := 0, 0; ip < dlen; ip, op = ip+4, op+1 {
		out[op] = color.RGBA{
			R: data[ip+0],
			G: data[ip+1],
			B: data[ip+2],
			A: data[ip+3],
		}
	}

	return out, nil
}

func convertRGBAToRGBA(img *image.RGBA) (out *Image, copied bool) {
	return convertRGBAToRGBASlow(img)
}
