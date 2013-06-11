package chunk

import (
	"math/big"
	"math/rand"
	"testing"
)

func TestGetSetUnset(t *testing.T) {
	zero := big.NewInt(0)
	c := New(zero, zero, zero)

	for x := uint8(0); x < 64; x++ {
		for y := uint8(0); y < 64; y++ {
			for z := uint8(0); z < 64; z++ {
				if c.Get(x, y, z) {
					t.Errorf("[%d %d %d] started set", x, y, z)
				}
			}
		}
	}

	r := rand.New(rand.NewSource(0))
	for x := uint8(0); x < 64; x++ {
		for y := uint8(0); y < 64; y++ {
			for z := uint8(0); z < 64; z++ {
				if r.Intn(2) == 0 {
					c.Set(x, y, z)
				}
				if r.Intn(2) == 0 {
					c.Unset(x, y, z)
				}
			}
		}
	}

	r = rand.New(rand.NewSource(0))
	for x := uint8(0); x < 64; x++ {
		for y := uint8(0); y < 64; y++ {
			for z := uint8(0); z < 64; z++ {
				expected := r.Intn(2) == 0
				expected = r.Intn(2) != 0 && expected
				if actual := c.Get(x, y, z); actual != expected {
					t.Errorf("[%d %d %d] expected:%v actual:%v", x, y, z, expected, actual)
				}
			}
		}
	}
}
