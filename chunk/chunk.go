package chunk

import "math/big"

type Chunk struct {
	X, Y, Z big.Int
	Solid   [64][64]uint64
}

// Constructs a new *Chunk that contains the global point (64x, 64y, 64z). The returned *Chunk
// is completely non-solid.
func New(x, y, z *big.Int) *Chunk {
	const shift = 6 // 1<<6 == 64

	c := &Chunk{}

	c.X.Set(x)
	c.Y.Set(y)
	c.Z.Set(z)

	return c
}

// Returns true if the local point (x, y, z) is solid. Coordinates are in the range [0,64).
// Getting the solidity of a point outside the legal range has undefined results.
func (c *Chunk) Get(x, y, z uint8) bool {
	return c.Solid[x][y]>>z&1 != 0
}

// Makes the local point (x, y, z) solid. Coordinates are in the range [0,64).
// Setting the solidity of a point outside the legal range has undefined results.
func (c *Chunk) Set(x, y, z uint8) {
	c.Solid[x][y] |= 1 << z
}

// Makes the local point (x, y, z) non-solid. Coordinates are in the range [0,64).
// Unsetting the solidity of a point outside the legal range has undefined results.
func (c *Chunk) Unset(x, y, z uint8) {
	c.Solid[x][y] &^= 1 << z
}

func (c *Chunk) MarshalJSON() ([]byte, error) {
	x, err := c.X.MarshalJSON()
	if err != nil {
		return nil, err
	}
	y, err := c.Y.MarshalJSON()
	if err != nil {
		return nil, err
	}
	z, err := c.Z.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var b []byte
	b = append(b, `{"X":`...)
	b = append(b, x...)
	b = append(b, `,"Y":`...)
	b = append(b, y...)
	b = append(b, `,"Z":`...)
	b = append(b, z...)
	b = append(b, `,"Solid":[`...)
	for x := range c.Solid {
		b = append(b, `[`...)
		for y := range c.Solid[x] {
			z := c.Solid[x][y]
			if z == 0 {
				b = append(b, `"0",`...)
			} else {
				const digits = `0123456789`
				var num [23]byte
				n := len(num)
				n--
				num[n] = ','
				n--
				num[n] = '"'
				for z > 0 {
					n--
					num[n] = digits[z%10]
					z /= 10
				}
				n--
				num[n] = '"'
				b = append(b, num[n:]...)
			}
		}
		b = append(b[:len(b)-1], `],`...)
	}
	b = append(b[:len(b)-1], `]}`...)
	return b, nil
}

func (c *Chunk) Equals(o *Chunk) bool {
	return c.X.Cmp(&o.X) == 0 &&
		c.Y.Cmp(&o.Y) == 0 &&
		c.Z.Cmp(&o.Z) == 0 &&
		c.Solid == o.Solid

}

func (c *Chunk) String() string {
	return "Chunk[" + c.X.String() + "," + c.Y.String() + "," + c.Z.String() + "]"
}
