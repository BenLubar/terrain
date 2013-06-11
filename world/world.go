package world

import (
	"math/big"

	"github.com/BenLubar/terrain/chunk"
)

type World struct {
	// Exported for encoders. Do not access directly.
	Chunks map[string]*chunk.Chunk
}

func New() *World {
	w := &World{}
	w.Chunks = make(map[string]*chunk.Chunk)
	return w
}

func split(gc *big.Int, cc *big.Int) uint8 {
	cc.Rsh(gc, 6)
	return uint8(gc.Uint64() & 63)
}

func (w *World) Get(x, y, z *big.Int) bool {
	var cx, cy, cz big.Int
	lx := split(x, &cx)
	ly := split(y, &cy)
	lz := split(z, &cz)
	key := cx.String() + "," + cy.String() + "," + cz.String()
	c, ok := w.Chunks[key]
	if ok {
		return c.Get(lx, ly, lz)
	}

	return z.Cmp(zero) <= 0
}

var zero = big.NewInt(0)
var bottomSolid [64][64]uint64
var allSolid [64][64]uint64

func init() {
	for i := 0; i < 64; i++ {
		allSolid[0][i] = ^uint64(0)
		bottomSolid[0][i] = uint64(1)
	}
	for i := 1; i < 64; i++ {
		allSolid[i] = allSolid[0]
		bottomSolid[i] = bottomSolid[0]
	}
}

func (w *World) Chunk(x, y, z *big.Int) *chunk.Chunk {
	key := x.String() + "," + y.String() + "," + z.String()
	c, ok := w.Chunks[key]
	if ok {
		return c
	}

	c = chunk.New(x, y, z)
	switch z.Cmp(zero) {
	case 0:
		c.Solid = bottomSolid
	case -1:
		c.Solid = allSolid
	}

	w.Chunks[key] = c
	return c
}
