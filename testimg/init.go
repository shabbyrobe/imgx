package testimg

import "math/rand"

var defaultRNG = rand.New(rand.NewSource(0))

var defaultRandPalette = RandPalette(defaultRNG, 256)
