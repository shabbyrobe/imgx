//+build purego

package rgba

import (
	"fmt"
	"image"
	"image/color"
)

func castFromBytes(data []byte) ([]color.RGBA, error) {
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

func castToBytes(colors []color.RGBA) (data []byte, err error) {
	data := make([]byte, len(colors)*4)

	for ip, op := 0, 0; ip < dlen; ip, op = ip+1, op+4 {
		data[op+0] = colors[ip].R
		data[op+1] = colors[ip].G
		data[op+2] = colors[ip].B
		data[op+3] = colors[ip].A
	}

	return data, nil
}

func convertRGBAToRGBA(img *image.RGBA) (out *Image, copied bool) {
	return convertRGBAToRGBASlow(img)
}
