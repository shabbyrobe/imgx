package rgba

import (
	"fmt"
	"image/color"
)

func NRGBAFromHex(s string) (c color.NRGBA, err error) {
	if s[0] == '#' {
		s = s[1:]
	}

	hasAlpha := len(s) == 8
	if len(s) != 6 && !hasAlpha {
		return c, fmt.Errorf("rgba: invalid hex %q", s)
	}

	r := (hexVals[s[0]] << 4) | hexVals[s[1]]
	g := (hexVals[s[2]] << 4) | hexVals[s[3]]
	b := (hexVals[s[4]] << 4) | hexVals[s[5]]
	if r < 0 || g < 0 || b < 0 {
		return c, fmt.Errorf("rgb: invalid hex %q", s)
	}
	a := int(0xff)
	if hasAlpha {
		a = (hexVals[s[6]] << 4) | hexVals[s[7]]
	}
	c = color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	return c, nil
}

func RGBAFromHex(s string) (c color.RGBA, err error) {
	nr, err := NRGBAFromHex(s)
	if err != nil {
		return c, err
	}
	a := uint32(nr.A)
	return color.RGBA{
		R: uint8(uint32(nr.R) * a / 0xff),
		G: uint8(uint32(nr.G) * a / 0xff),
		B: uint8(uint32(nr.B) * a / 0xff),
		A: uint8(a),
	}, nil
}

var (
	hexVals = [256]int{}
)

func init() {
	for i := 0; i < 256; i++ {
		// This will always produce a negative result in our conversion, even if doing
		// RGBA64s, i.e: '(hexVals[b[0]] << 8) | (hexVals[b[1]])' is guaranteed negative
		// if either b[0] or b[1] is invalid
		hexVals[i] = -0x1_0000
	}
	for i := 0; i < 10; i++ {
		hexVals['0'+i] = i
	}
	for i := 0; i < 6; i++ {
		hexVals['a'+i] = 10 + i
		hexVals['A'+i] = 10 + i
	}
}
