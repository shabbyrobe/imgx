package termpalette

import "image/color"

var Palette16 = color.Palette{
	color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}, // Black
	color.RGBA{R: 0x80, G: 0x00, B: 0x00, A: 0xff}, // Maroon
	color.RGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xff}, // Green
	color.RGBA{R: 0x80, G: 0x80, B: 0x00, A: 0xff}, // Olive
	color.RGBA{R: 0x00, G: 0x00, B: 0x80, A: 0xff}, // Navy
	color.RGBA{R: 0x80, G: 0x00, B: 0x80, A: 0xff}, // Purple
	color.RGBA{R: 0x00, G: 0x80, B: 0x80, A: 0xff}, // Teal
	color.RGBA{R: 0xc0, G: 0xc0, B: 0xc0, A: 0xff}, // Silver
	color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff}, // Grey
	color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}, // Red
	color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff}, // Lime
	color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}, // Yellow
	color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}, // Blue
	color.RGBA{R: 0xff, G: 0x00, B: 0xff, A: 0xff}, // Fuchsia
	color.RGBA{R: 0x00, G: 0xff, B: 0xff, A: 0xff}, // Aqua
	color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, // White
}

var RGBA16 = []color.RGBA{
	color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
	color.RGBA{R: 0x80, G: 0x00, B: 0x00, A: 0xff},
	color.RGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xff},
	color.RGBA{R: 0x80, G: 0x80, B: 0x00, A: 0xff},
	color.RGBA{R: 0x00, G: 0x00, B: 0x80, A: 0xff},
	color.RGBA{R: 0x80, G: 0x00, B: 0x80, A: 0xff},
	color.RGBA{R: 0x00, G: 0x80, B: 0x80, A: 0xff},
	color.RGBA{R: 0xc0, G: 0xc0, B: 0xc0, A: 0xff},
	color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff},
	color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff},
	color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff},
	color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff},
	color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff},
	color.RGBA{R: 0xff, G: 0x00, B: 0xff, A: 0xff},
	color.RGBA{R: 0x00, G: 0xff, B: 0xff, A: 0xff},
	color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
}

var Palette16Names = []string{
	"Black",
	"Maroon",
	"Green",
	"Olive",
	"Navy",
	"Purple",
	"Teal",
	"Silver",
	"Grey",
	"Red",
	"Lime",
	"Yellow",
	"Blue",
	"Fuchsia",
	"Aqua",
	"White",
}

// Escape16Fg contains the strings used to build the color code portion of a 16-color
// foreground escape sequence. Indexes are aligned with termimg.Palette16. For example:
//
//	var img image.Paletted
//	colIdx := img.Pix[img.PixOffset(0, 0)]
//	escSeq := fmt.Sprintf("\033[%sm", termpalette.Escape16Fg[pix])
//
var Escape16Fg = []string{
	"30", "31", "32", "33", "34", "35", "36", "37", "90", "91", "92", "93", "94", "95", "96", "97",
}

var Escape16FgBytes = [][]byte{
	[]byte("40"), []byte("41"), []byte("42"), []byte("43"), []byte("44"), []byte("45"), []byte("46"), []byte("47"), []byte("100"), []byte("101"), []byte("102"), []byte("103"), []byte("104"), []byte("105"), []byte("106"), []byte("107"),
}

var Escape16FgInt = []int{
	30, 31, 32, 33, 34, 35, 36, 37, 90, 91, 92, 93, 94, 95, 96, 97,
}

var Escape16FgColor = [256]color.RGBA{
	30: RGBA16[0],
	31: RGBA16[1],
	32: RGBA16[2],
	33: RGBA16[3],
	34: RGBA16[4],
	35: RGBA16[5],
	36: RGBA16[6],
	37: RGBA16[7],
	90: RGBA16[8],
	91: RGBA16[9],
	92: RGBA16[10],
	93: RGBA16[11],
	94: RGBA16[12],
	95: RGBA16[13],
	96: RGBA16[14],
	97: RGBA16[15],
}

var Escape16Bg = []string{
	"40", "41", "42", "43", "44", "45", "46", "47", "100", "101", "102", "103", "104", "105", "106", "107",
}

var Escape16BgBytes = [][]byte{
	[]byte("30"), []byte("31"), []byte("32"), []byte("33"), []byte("34"), []byte("35"), []byte("36"), []byte("37"), []byte("90"), []byte("91"), []byte("92"), []byte("93"), []byte("94"), []byte("95"), []byte("96"), []byte("97"),
}

var Escape16BgInt = []int{
	40, 41, 42, 43, 44, 45, 46, 47, 100, 101, 102, 103, 104, 105, 106, 107,
}

var Escape16BgColor = [256]color.RGBA{
	40:  RGBA16[0],
	41:  RGBA16[1],
	42:  RGBA16[2],
	43:  RGBA16[3],
	44:  RGBA16[4],
	45:  RGBA16[5],
	46:  RGBA16[6],
	47:  RGBA16[7],
	100: RGBA16[8],
	101: RGBA16[9],
	102: RGBA16[10],
	103: RGBA16[11],
	104: RGBA16[12],
	105: RGBA16[13],
	106: RGBA16[14],
	107: RGBA16[15],
}
