package termpalette

// Code contains the strings used to build the color code portion of a 256 color
// escape sequence. Indexes are aligned with termimg.Palette. For example:
//
//	var img image.Paletted
//	pix := img.Pix[img.PixOffset(0, 0)]
//	esc := fmt.Sprintf("\033[38;5;%sm", termpalette.Code[pix])
//
// Code is the string equivalent of CodeBytes.
//
var Code = codeString

// CodeBytes contains the bytes used to build the color code portion of a 256 color
// escape sequence. It is the []byte equivalent of Code. See Code for more info.
//
var CodeBytes = codeBytes

var CodeInt = codeInt
