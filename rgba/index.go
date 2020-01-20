package rgba

import "image/color"

// Indexer is used to create an Index of a Palette. The Index allows
// lookups of nearest-neighbour colour values.
//
// See NewRGBAPrecacheIndexer, NewRGBPrecacheIndexer and NewRGBATreeIndexer.
//
// Also see IndexUnmarshaler.
type Indexer interface {
	IndexRGBAPalette(pal Palette) Index
}

// Index is used to perform lookups of nearest-neighbour colour values.
// It is used to implement faster versions of image/color.Palette.Convert.
//
// See NewRGBAPrecacheIndexer, NewRGBPrecacheIndexer and NewRGBATreeIndexer.
type Index interface {
	NearestRGBA(c color.RGBA) (nn color.RGBA, idx int)
	NearestRGBAIndex(c color.RGBA) int
	NearestRGBAColor(c color.RGBA) color.RGBA
}

// IndexUnmarshaler may be implemented by an Indexer to allow unserializing
// a binary representation of an Index.
//
// This can be handy when the index might be slower to construct than desired,
// or when a precomputed index makes sense to bake into your binary.
//
// None of the marshaled formats are guaranteed to be stable. It is highly recommended
// that they only be used in-process, or in a code-generation step (which is run any time
// go.mod updates).
//
// See RGBPrecacheIndexer.
//
type IndexUnmarshaler interface {
	UnmarshalIndex(b []byte) (Index, error)
}

// IndexMarshaler may be implemented by an Index to allow serializing a binary
// representation for faster Unmarshaling later.
//
// See RGBPrecacheIndex.
//
type IndexMarshaler interface {
	MarshalIndex() []byte
}
