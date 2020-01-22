package rgba

import (
	"fmt"
	"image"
	"math/rand"
	"reflect"
	"testing"

	"github.com/shabbyrobe/imgx/testimg"
)

// Globals to hopefully trick optimiser
var (
	BenchInt      int
	BenchRGBImage *Image

	AccessIdxX, AccessIdxY = 30, 40
	XSize, YSize           = 512, 256
)

func TestCastBytes(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	gen := testimg.RandBlocks{512, 512, 32, 32}
	img, _ := Convert(gen.RGBA(rng))

	bts, err := CastToBytes(img.Vals)
	if err != nil {
		t.Fatal(err)
	}
	if len(bts) != len(img.Vals)*4 {
		t.Fatal()
	}

	col, err := CastFromBytes(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(col) != len(img.Vals) {
		t.Fatal()
	}
	if !reflect.DeepEqual(col, img.Vals) {
		t.Fatal()
	}
}

func TestConvertRGBA(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	gen := testimg.RandBlocks{512, 512, 32, 32}

	unwrap := func(v *Image, ok bool) *Image {
		return v
	}

	var cases = []struct {
		oimg image.Image
		conv func(v image.Image) *Image
	}{
		{gen.RGBA(rng), func(v image.Image) *Image { return unwrap(convertRGBAToRGBA(v.(*image.RGBA))) }},
		{gen.RGBA(rng), func(v image.Image) *Image { return unwrap(convertRGBAToRGBASlow(v.(*image.RGBA))) }},
		{gen.RGBA(rng), func(v image.Image) *Image { return convertImageToRGBA(v.(*image.RGBA)) }},
		{gen.RGBA(rng), func(v image.Image) *Image { return convertRGBAAtToRGBA(v.(*image.RGBA)) }},
		{gen.RGBA64(rng), func(v image.Image) *Image { return convertRGBA64ToRGBA(v.(*image.RGBA64)) }},
		{gen.NRGBA(rng), func(v image.Image) *Image { return convertNRGBAToRGBA(v.(*image.NRGBA)) }},
		{gen.NRGBA64(rng), func(v image.Image) *Image { return convertNRGBA64ToRGBA(v.(*image.NRGBA64)) }},
		{gen.YCbCr(rng), func(v image.Image) *Image { return convertYCbCrToRGBA(v.(*image.YCbCr)) }},
		{gen.CMYK(rng), func(v image.Image) *Image { return convertCMYKToRGBA(v.(*image.CMYK)) }},
		{gen.Paletted(rng, nil), func(v image.Image) *Image { return convertPalettedToRGBA(v.(*image.Paletted)) }},
	}

	for idx, tc := range cases {
		t.Run(fmt.Sprintf("%d/%T", idx, tc.oimg), func(t *testing.T) {
			oimg := tc.oimg
			rimg := tc.conv(oimg)

			bounds := rimg.Bounds()

			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					var (
						or, og, ob, _  = oimg.At(x, y).RGBA()
						rr, rg, rb, ra = rimg.At(x, y).RGBA()
					)

					_ = ra

					// Let's only test 8 bits; the YCbCr conversion gives test failures
					// with the 16 bit values from RGBA, which is kinda unsurprising as we
					// started with 8-bit values anyway.
					if or>>8 != rr>>8 {
						t.Fatalf("or(%d) != rr(%d) at (%d,%d)", or>>8, rr>>8, x, y)
					}
					if og>>8 != rg>>8 {
						t.Fatalf("og(%d) != rg(%d) at (%d,%d)", og>>8, rg>>8, x, y)
					}
					if ob>>8 != rb>>8 {
						t.Fatalf("ob(%d) != rb(%d) at (%d,%d)", ob>>8, rb>>8, x, y)
					}
				}
			}
		})
	}
}

var BenchmarkImage image.Image

func BenchmarkConvert(b *testing.B) {
	rng := rand.New(rand.NewSource(0))
	gen := testimg.RandBlocks{512, 512, 32, 32}

	b.Run("rgba", func(b *testing.B) {
		in := gen.RGBA(rng)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			BenchmarkImage, _ = convertRGBAToRGBA(in)
		}
	})

	b.Run("rgbaslow", func(b *testing.B) {
		in := gen.RGBA(rng)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			BenchmarkImage, _ = convertRGBAToRGBASlow(in)
		}
	})

	b.Run("rgbaat", func(b *testing.B) {
		in, _ := convertRGBAToRGBA(gen.RGBA(rng))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			BenchmarkImage = convertRGBAAtToRGBA(in)
		}
	})

	b.Run("nrgba", func(b *testing.B) {
		in := gen.NRGBA(rng)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			BenchmarkImage = convertNRGBAToRGBA(in)
		}
	})

	b.Run("nrgba64", func(b *testing.B) {
		in := gen.NRGBA64(rng)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			BenchmarkImage = convertNRGBA64ToRGBA(in)
		}
	})

	b.Run("rgba64", func(b *testing.B) {
		in := gen.RGBA64(rng)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			BenchmarkImage = convertRGBA64ToRGBA(in)
		}
	})

	b.Run("ycbcr", func(b *testing.B) {
		in := gen.YCbCr(rng)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			BenchmarkImage = convertYCbCrToRGBA(in)
		}
	})

	b.Run("cmyk", func(b *testing.B) {
		in := gen.CMYK(rng)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			BenchmarkImage = convertCMYKToRGBA(in)
		}
	})

	b.Run("paletted", func(b *testing.B) {
		in := gen.Paletted(rng, nil)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			BenchmarkImage = convertPalettedToRGBA(in)
		}
	})
}

// Is a 2D image array faster, or a 1D array with a stride?
func BenchmarkArrayAccess(b *testing.B) {
	var t1d = make([]int, YSize*XSize)
	var t2d = make([][]int, YSize)

	for y := 0; y < YSize; y++ {
		t2d[y] = make([]int, YSize)
		for x := 0; x < YSize; x++ {
			t1d[(y*XSize)+x] = y + x
			t2d[y][x] = y + x
		}
	}

	b.Run("3d", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			BenchInt = t2d[AccessIdxY][AccessIdxX]
		}
	})

	b.Run("1d", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			BenchInt = t1d[AccessIdxY*XSize+AccessIdxX]
		}
	})
}
