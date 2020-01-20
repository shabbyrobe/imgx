package rgba

import (
	"encoding/binary"
	"fmt"
	"image/color"
)

// NewRGBPrecacheIndexer creates an indexer which can build slightly inaccurate
// pre-computed mappings of 5-bit RGB values to their nearest neighbour in the provided
// palette, discarding alpha.
//
// It is at least an order of magnitude faster than NewRGBTreeIndexer and I honestly
// can't tell the difference just by looking with my own eyes. YMMV.
//
// The colors returned by the resulting indexes may or may not have their alpha
// replaced with 0xff - the current behaviour should not be relied upon.
//
// There is no RGBAPrecacheIndexer; the brute force code to build the index was
// too slow and the Index used too much memory.
//
func NewRGBPrecacheIndexer(using Indexer) Indexer {
	if using == nil {
		using = NewRGBTreeIndexer()
	}
	return &rgbPrecacheIndexer{using: using}
}

type rgbPrecacheIndexer struct {
	using Indexer
}

func (pc rgbPrecacheIndexer) IndexRGBAPalette(pal Palette) Index {
	ix := pc.using.IndexRGBAPalette(pal)
	rp := &rgbPrecacheIndex{pal: pal}

	for r := 0; r < 32; r++ {
		for g := 0; g < 32; g++ {
			for b := 0; b < 32; b++ {
				col := color.RGBA{
					R: uint8(r << 3),
					G: uint8(g << 3),
					B: uint8(b << 3),
					A: 0xff,
				}

				ic, ii := ix.NearestRGBA(col)
				rp.color[r][g][b] = color.RGBA{ic.R, ic.G, ic.B, 0xff}
				rp.index[r][g][b] = int32(ii)
			}
		}
	}

	return rp
}

type rgbPrecacheIndex struct {
	pal   Palette
	color [32][32][32]color.RGBA
	index [32][32][32]int32
}

func (pc *rgbPrecacheIndex) NearestRGBAIndex(c color.RGBA) int {
	return int(pc.index[c.R>>3][c.G>>3][c.B>>3])
}

func (pc *rgbPrecacheIndex) NearestRGBAColor(c color.RGBA) color.RGBA {
	return pc.color[c.R>>3][c.G>>3][c.B>>3]
}

func (pc *rgbPrecacheIndex) NearestRGBA(c color.RGBA) (nn color.RGBA, idx int) {
	nn = pc.color[c.R>>3][c.G>>3][c.B>>3]
	idx = int(pc.index[c.R>>3][c.G>>3][c.B>>3])
	return nn, idx
}

func (rgbPrecacheIndexer) UnmarshalIndex(data []byte) (Index, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("rgba: invalid data size")
	}

	var pc rgbPrecacheIndex
	palsz := int(binary.LittleEndian.Uint32(data))
	pc.pal = make(Palette, palsz)
	pos := 4

	for i := 0; i < palsz; i++ {
		pc.pal[i] = color.RGBA{
			R: data[pos],
			G: data[pos+1],
			B: data[pos+2],
			A: data[pos+3],
		}
		pos += 4
	}

	idx := int64(0)
	for c := 0; c < 32768; {
		run, n := binary.Varint(data[pos:])
		if n <= 0 || run <= 0 {
			return nil, fmt.Errorf("rgba: invalid index")
		}
		pos += n

		delt, n := binary.Varint(data[pos:])
		if n <= 0 {
			return nil, fmt.Errorf("rgba: invalid index")
		}
		pos += n

		idx = int64(idx + delt)

		for r := int64(0); r < run; r, c = r+1, c+1 {
			if int(idx) >= palsz {
				return nil, fmt.Errorf("rgba: invalid palette index")
			}

			r, g, b := c>>10, (c>>5)&0b11111, c&0b11111
			pc.index[r][g][b] = int32(idx)
			pc.color[r][g][b] = pc.pal[idx]
		}
	}

	return &pc, nil
}

func (pc *rgbPrecacheIndex) MarshalIndex() []byte {
	var bts = make([]byte, 0, 65536)

	var scratchArr [32]byte
	var scratch = scratchArr[:]

	binary.LittleEndian.PutUint32(scratch, uint32(len(pc.pal)))
	bts = append(bts, scratch[:4]...)

	for _, col := range pc.pal {
		bts = append(bts, col.R, col.G, col.B, col.A)
	}

	var last, lastDelt, run int64
	for c := 0; c < 32768; c++ {
		v := pc.index[c>>10][(c>>5)&0b11111][c&0b11111]
		delt := int64(v) - last
		last = int64(v)

		if delt != 0 {
			if run > 0 {
				n := binary.PutVarint(scratch, run)
				bts = append(bts, scratch[:n]...)
				n = binary.PutVarint(scratch, lastDelt)
				bts = append(bts, scratch[:n]...)
			}
			lastDelt, run = delt, 1
		} else {
			run++
		}
	}

	n := binary.PutVarint(scratch, run)
	bts = append(bts, scratch[:n]...)
	n = binary.PutVarint(scratch, lastDelt)
	bts = append(bts, scratch[:n]...)

	return bts
}
