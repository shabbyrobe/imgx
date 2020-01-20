package termpalette

// Code contains the strings used to build the color code portion of a 256 color
// escape sequence. Indexes are aligned with termimg.Palette. For example:
//
//	var img image.Paletted
//	colIdx := img.Pix[img.PixOffset(0, 0)]
//	escSeq := fmt.Sprintf("\033[38;5;%sm", termpalette.Code[pix])
//
var Code = codeString

// Code contains the string used to build the color code portion of a 256 color
// escape sequence. See Code for more info.
//
var CodeBytes = codeBytes

var CodeInt = codeInt
