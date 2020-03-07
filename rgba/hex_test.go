package rgba

import (
	"math/rand"
	"testing"

	"github.com/shabbyrobe/imgx/testimg"
)

func TestRGBAHexRand(t *testing.T) {
	rng := rand.New(rand.NewSource(0))

	for i := 0; i < 10000; i++ {
		col := testimg.RandRGBA(rng)

		hex := ToNRGBAHex(col)
		back, err := FromNRGBAHex(hex)
		if err != nil {
			t.Fatal(err)
		}

		diffR := int(back.R) - int(col.R)
		diffG := int(back.G) - int(col.G)
		diffB := int(back.B) - int(col.B)
		diffA := int(back.A) - int(col.A)

		// Unfortunately, 8-bit NRGBA to RGBA is not quite accurate.
		const limit = 1
		if diffR < -limit || diffR > limit {
			t.Fatalf("back %#v != col %#v", back, col)
		}
		if diffG < -limit || diffG > limit {
			t.Fatalf("back %#v != col %#v", back, col)
		}
		if diffB < -limit || diffB > limit {
			t.Fatalf("back %#v != col %#v", back, col)
		}
		if diffA < -limit || diffA > limit {
			t.Fatalf("back %#v != col %#v", back, col)
		}
	}
}
