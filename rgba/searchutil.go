package rgba

func sqDiff8(x, y uint8) uint32 {
	d := uint16(x) - uint16(y)
	return uint32(d * d) // uint32 allows us to add 4 of these without overflow
}
