package rgba

import (
	"fmt"
	"image/color"
	"math"
	"sort"
)

func NewRGBTreeIndexer() Indexer {
	return &rgbTreeIndexer{}
}

type rgbTreeIndexer struct{}

func (rgbTreeIndexer) IndexRGBAPalette(pal Palette) Index {
	return rgbTreeBuild(pal)
}

type rgbNode struct {
	left  *rgbNode
	right *rgbNode
	index int
	axis  rgbAxis
	col   color.RGBA
}

func (kd rgbNode) RGBA() color.RGBA {
	return kd.col
}

func (kd *rgbNode) NearestRGBAIndex(c color.RGBA) int {
	node, _ := kd.nnRecursive(c, math.MaxUint32, 0)
	return node.index
}

func (kd *rgbNode) NearestRGBAColor(c color.RGBA) color.RGBA {
	node, _ := kd.nnRecursive(c, math.MaxUint32, 0)
	return node.col
}

func (kd *rgbNode) NearestRGBA(c color.RGBA) (col color.RGBA, idx int) {
	node, _ := kd.nnRecursive(c, math.MaxUint32, 0)
	return node.col, node.index
}

func (kd *rgbNode) nnRecursive(c color.RGBA, best uint32, depth int) (nn *rgbNode, bestResult uint32) {
	// XXX(bw): I tried an iterative approach to this (rather than a recursive one) but it
	// was about the same speed, and vastly more horrible.

	nn = kd

	dist := sqDiff8(kd.col.R, c.R) + sqDiff8(kd.col.G, c.G) + sqDiff8(kd.col.B, c.B)
	if dist < best {
		best = dist
	}

	var nval, qval uint8

	switch kd.axis {
	case rgbAxisR:
		nval, qval = kd.col.R, c.R
	case rgbAxisG:
		nval, qval = kd.col.G, c.G
	case rgbAxisB:
		nval, qval = kd.col.B, c.B
	default:
		panic("unknown axis")
	}

	var cmp = int(qval) - int(nval)

	near, far := kd.left, kd.right
	if cmp > 0 {
		near, far = kd.right, kd.left
	}

	if near != nil {
		nextNode, nextBest := near.nnRecursive(c, best, depth+1)
		if nextBest < best {
			nn, best = nextNode, nextBest
		}
	}

	if far != nil {
		planeDist := uint32(cmp * cmp)

		if planeDist < best {
			farLeaf, farBest := far.nnRecursive(c, best, depth+1)
			if farBest < best {
				nn, best = farLeaf, farBest
			}
		}
	}

	return nn, best
}

type rgbTreeItem struct {
	index int
	col   color.RGBA
}

type rgbTreeBuilder struct {
	vals [256]int
	slab []rgbNode
	next int
}

func rgbTreeBuild(items []color.RGBA) *rgbNode {
	if len(items) > 256 {
		// FIXME: would be good to remove this limitation
		panic(fmt.Errorf("palette length must be <= 256"))
	}

	ilen := len(items)

	var bld = rgbTreeBuilder{
		slab: make([]rgbNode, ilen),
	}

	var bItems = make([]rgbTreeItem, len(items))
	for idx, col := range items {
		item := rgbTreeItem{col: col, index: idx}
		item.col.A = 0xff // Discard alpha!
		bItems[idx] = item
	}
	return bld.node(bItems, 0)
}

func (bld *rgbTreeBuilder) node(items []rgbTreeItem, axis rgbAxis) *rgbNode {
	node := &bld.slab[bld.next]
	node.axis = axis
	bld.next++

	switch len(items) {
	case 0:
		return nil
	case 1:
		node.col, node.index = items[0].col, items[0].index
		return node
	}

	nums := bld.vals[:0]
	for idx, item := range items {
		var v uint8
		switch axis {
		case rgbAxisR:
			v = item.col.R
		case rgbAxisG:
			v = item.col.G
		case rgbAxisB:
			v = item.col.B
		default:
			panic("unknown axis")
		}
		nums = append(nums, int(v)<<8|idx)
	}

	// FIXME: There might be a better option for this, go's general purpose sort is useful
	// but maybe we can do better as this accounts for 80% of the time this takes. A
	// sorting network doesn't work here for all cases (256 items is WAY too much
	// generated code and the compiler takes a billion years)
	if !networkSortInt(nums, len(nums)) {
		sort.Ints(nums)
	}

	// FIXME: we can do something fancy later where we allocate once for all depths
	// (for a 32 color palette, we can allocate 32+16+8+4+2+1 and use the depth
	// to pick which power of 2 segment of sortedItems is safe to use without clobbering
	// parts being used by other levels of the stack)
	sortedItems := make([]rgbTreeItem, len(nums))
	for i, n := range nums {
		sortedItems[i] = items[n&0xff]
	}

	// do not use the 'items' or 'nums' vars below this point

	medianIndex := len(nums) / 2
	median := &sortedItems[medianIndex]
	node.col, node.index = median.col, median.index

	leftItems := sortedItems[:medianIndex]
	rightItems := sortedItems[medianIndex+1:]
	if len(leftItems) != 0 {
		node.left = bld.node(leftItems, axis.Next())
	}
	if len(rightItems) != 0 {
		node.right = bld.node(rightItems, axis.Next())
	}
	return node
}

type rgbAxis int

func (x rgbAxis) Next() rgbAxis {
	if x != rgbAxisB {
		return x + 1
	} else {
		return rgbAxisR
	}
}

func (x rgbAxis) String() string {
	var v string
	switch x {
	case rgbAxisR:
		v = "R"
	case rgbAxisG:
		v = "G"
	case rgbAxisB:
		v = "B"
	default:
		panic("unknown axis")
	}
	return fmt.Sprintf("%s(%d)", v, x)
}

const (
	rgbAxisR rgbAxis = iota
	rgbAxisG
	rgbAxisB
)
