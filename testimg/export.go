//+build ignore

package main

import (
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/shabbyrobe/termimg/internal/termpalette"
	"github.com/shabbyrobe/termimg/internal/testimg"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func encodePNGFile(img image.Image, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	return png.Encode(f, img)
}

func run() error {
	dir := os.Args[1]

	for _, c := range []struct {
		File   string
		Recipe testimg.Recipe
	}{
		{"testimg-123x45-1x1", testimg.RandBlocks{123, 45, 1, 1}},
		{"testimg-512x512-1x1", testimg.RandBlocks{512, 512, 1, 1}},
		{"testimg-512x512-512x1", testimg.RandBlocks{512, 512, 512, 1}},
		{"testimg-512x512-1x512", testimg.RandBlocks{512, 512, 1, 512}},
		{"testimg-512x512-2x2", testimg.RandBlocks{512, 512, 2, 2}},
		{"testimg-512x512-8x8", testimg.RandBlocks{512, 512, 8, 8}},
		{"testimg-512x512-32x32", testimg.RandBlocks{512, 512, 32, 32}},
	} {
		fn := filepath.Join(dir, c.File)

		rng := rand.New(rand.NewSource(0))
		if err := encodePNGFile(c.Recipe.RGBA(rng), fn+"-24bit.png"); err != nil {
			return err
		}

		rng = rand.New(rand.NewSource(0))
		if err := encodePNGFile(c.Recipe.Paletted(rng, termpalette.Palette), fn+"-term.png"); err != nil {
			return err
		}

		rng = rand.New(rand.NewSource(0))
		if err := encodePNGFile(c.Recipe.YCbCr(rng), fn+"-ycbcr.png"); err != nil {
			return err
		}

		rng = rand.New(rand.NewSource(0))
		if err := encodePNGFile(c.Recipe.CMYK(rng), fn+"-cmyk.png"); err != nil {
			return err
		}

		rng = rand.New(rand.NewSource(0))
		if err := encodePNGFile(c.Recipe.RGBA64(rng), fn+"-rgba64.png"); err != nil {
			return err
		}

		rng = rand.New(rand.NewSource(0))
		if err := encodePNGFile(c.Recipe.NRGBA64(rng), fn+"-nrgba64.png"); err != nil {
			return err
		}
	}

	return nil
}
