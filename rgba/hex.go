package rgba

import (
	"fmt"
	"image/color"
)

func NRGBAFromNRGBAHex(s string) (c color.NRGBA, err error) {
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

func FromNRGBAHex(s string) (c color.RGBA, err error) {
	nr, err := NRGBAFromNRGBAHex(s)
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

// If len(b) < 9, this will panic.
func WriteNRGBAHex(c color.RGBA, buf []byte) {
	_ = buf[8]

	a := uint32(c.A)
	if a == 0 {
		copy(buf, "#00000000")
		return
	}
	if a == 0xff {
		buf[0] = '#'
		copy(buf[1:3], hexStrs[int(c.R)+int(c.R):])
		copy(buf[3:5], hexStrs[int(c.G)+int(c.G):])
		copy(buf[5:7], hexStrs[int(c.B)+int(c.B):])
		buf[7] = 'F'
		buf[8] = 'F'
		return
	}

	buf[0] = '#'
	r, g, b := uint32(c.R), uint32(c.G), uint32(c.B)
	r = (r * 0xff) / a
	g = (g * 0xff) / a
	b = (b * 0xff) / a
	copy(buf[1:3], hexStrs[r+r:])
	copy(buf[3:5], hexStrs[g+g:])
	copy(buf[5:7], hexStrs[b+b:])
	copy(buf[7:9], hexStrs[a+a:])
}

func ToNRGBAHex(c color.RGBA) string {
	var buf = make([]byte, 9)
	WriteNRGBAHex(c, buf)
	return string(buf)
}

var (
	hexVals = [256]int{}
)

const hexStrs = "000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2D2E2F303132333435363738393A3B3C3D3E3F404142434445464748494A4B4C4D4E4F505152535455565758595A5B5C5D5E5F606162636465666768696A6B6C6D6E6F707172737475767778797A7B7C7D7E7F808182838485868788898A8B8C8D8E8F909192939495969798999A9B9C9D9E9FA0A1A2A3A4A5A6A7A8A9AAABACADAEAFB0B1B2B3B4B5B6B7B8B9BABBBCBDBEBFC0C1C2C3C4C5C6C7C8C9CACBCCCDCECFD0D1D2D3D4D5D6D7D8D9DADBDCDDDEDFE0E1E2E3E4E5E6E7E8E9EAEBECEDEEEFF0F1F2F3F4F5F6F7F8F9FAFBFCFDFEFF"

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
