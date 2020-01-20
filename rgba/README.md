rgba - Alternative image format to image.RGBA
=============================================

`rgba` implements an alternative image format to `image.RGBA`. I feel `rgba`
produces more readable code with less array-index gymnastics. I use this as the
basis for a lot of further processing.

I felt it was a mistake for a while, so I tried to migrate my code to use
`image.RGBA` and ran screaming back to the warm embrace of `rgba`.

Converting from an `image.RGBA` to an `rgba.Image` should be basically
"cheap-as-free" on targets that support `unsafe` code. See `rgba.Convert()`.

This package also contains a fast, "accurate"* colour search implementation
(see `NewRGBATreeIndexer()` and `NewRGBTreeIndexer()`) based on a kd-tree, and
a very fast, slightly inaccurate colour search implementation (see
`NewRGBPrecacheIndexer()`) based on a precomputed cache of all possible RGB
values after lopping off the three least significant bits. You probably want to
use the Precache indexer unless you need alpha or high accuracy. I personally
can't tell the difference visually.

* "accurate", insofar as one can be when pretending RGB -- whatever that means --
  is a cartesian space with orthogonal axes!

