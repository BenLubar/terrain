package world

import (
	"math/big"

	"github.com/petar/GoLLRB/llrb"

	"github.com/BenLubar/terrain/chunk"
)

type World struct {
	chunks *llrb.LLRB
}

type item struct {
	X, Y, Z *big.Int
	*chunk.Chunk
}

func (i *item) Less(than llrb.Item) bool {
	o := than.(*item)
	switch i.X.Cmp(o.X) {
	case -1:
		return true
	case 1:
		return false
	}
	switch i.Y.Cmp(o.Y) {
	case -1:
		return true
	case 1:
		return false
	}
	switch i.Z.Cmp(o.Z) {
	case -1:
		return true
	case 1:
		return false
	}
	return false
}

// Constructs a new, empty world. All locations with z < 0 are solid. All other
// locations are non-solid.
func New() *World {
	w := &World{}
	w.chunks = llrb.New()
	return w
}

func split(gc *big.Int, cc *big.Int) uint8 {
	cc.Rsh(gc, 6)
	return uint8(gc.Uint64() & 63)
}

func (w *World) Get(x, y, z *big.Int) bool {
	i := &item{
		X: &big.Int{},
		Y: &big.Int{},
		Z: &big.Int{},
	}
	lx := split(x, i.X)
	ly := split(y, i.Y)
	lz := split(z, i.Z)
	i, _ = w.chunks.Get(i).(*item)
	if i != nil {
		return i.Get(lx, ly, lz)
	}

	return z.Cmp(zero) < 0
}

func (w *World) Set(x, y, z *big.Int) {
	var cx, cy, cz big.Int
	lx := split(x, &cx)
	ly := split(y, &cy)
	lz := split(z, &cz)

	w.Chunk(&cx, &cy, &cz).Set(lx, ly, lz)
}

func (w *World) Unset(x, y, z *big.Int) {
	var cx, cy, cz big.Int
	lx := split(x, &cx)
	ly := split(y, &cy)
	lz := split(z, &cz)

	w.Chunk(&cx, &cy, &cz).Unset(lx, ly, lz)
}

var zero = big.NewInt(0)
var allSolid [64][64]uint64

func init() {
	for i := 0; i < 64; i++ {
		allSolid[0][i] = ^uint64(0)
	}
	for i := 1; i < 64; i++ {
		allSolid[i] = allSolid[0]
	}
}

func (w *World) Chunk(x, y, z *big.Int) *chunk.Chunk {
	i := &item{
		X: x,
		Y: y,
		Z: z,
	}
	i, _ = w.chunks.Get(i).(*item)
	if i != nil {
		return i.Chunk
	}

	i = &item{
		X: (&big.Int{}).Set(x),
		Y: (&big.Int{}).Set(y),
		Z: (&big.Int{}).Set(z),

		Chunk: chunk.New(x, y, z),
	}
	if z.Cmp(zero) < 0 {
		i.Solid = allSolid
	}

	w.chunks.InsertNoReplace(i)
	return i.Chunk
}
