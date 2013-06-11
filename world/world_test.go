package world

import (
	"math/big"
	"testing"
)

func TestChunk(t *testing.T) {
	var (
		zero    = big.NewInt(0)
		one     = big.NewInt(1)
		verylow = big.NewInt(-1 << 63)
		w       = New()
	)

	c0 := w.Chunk(zero, zero, zero)
	if c0 == nil {
		t.Fatalf("World.Chunk returned nil!")
	}
	if c0.X.Cmp(zero) != 0 || c0.Y.Cmp(zero) != 0 || c0.Z.Cmp(zero) != 0 {
		t.Errorf("Chunk(0, 0, 0) returned %v", c0)
	}

	zero = big.NewInt(0) // make zero not reference equivelent to the original zero
	c1 := w.Chunk(zero, zero, zero)
	if c1 == nil {
		t.Fatalf("World.Chunk returned nil!")
	}
	if c0 != c1 {
		t.Errorf("World.Chunk returned different values for the same input!")
	}
	if c1.X.Cmp(zero) != 0 || c1.Y.Cmp(zero) != 0 || c1.Z.Cmp(zero) != 0 {
		t.Errorf("Chunk(0, 0, 0) returned %v", c1)
	}

	c2 := w.Chunk(zero, zero, one)
	if c2 == nil {
		t.Fatalf("World.Chunk returned nil!")
	}
	if c0 == c2 || c1 == c2 {
		t.Errorf("World.Chunk returned the same value for different input!")
	}
	if c2.X.Cmp(zero) != 0 || c2.Y.Cmp(zero) != 0 || c2.Z.Cmp(one) != 0 {
		t.Errorf("Chunk(0, 0, 1) returned %v", c2)
	}

	c3 := w.Chunk(zero, zero, verylow)
	if c3 == nil {
		t.Fatalf("World.Chunk returned nil!")
	}
	if c0 == c3 || c1 == c3 || c2 == c3 {
		t.Errorf("World.Chunk returned the same value for different input!")
	}
	if c3.X.Cmp(zero) != 0 || c3.Y.Cmp(zero) != 0 || c3.Z.Cmp(verylow) != 0 {
		t.Errorf("Chunk(0, 0, %v) returned %v", verylow, c3)
	}

	if len(w.Chunks) != 3 {
		t.Errorf("Expected 3 chunks, but World has %d", len(w.Chunks))
	}
}

func TestGet(t *testing.T) {
	w := New()
	w.Chunk(big.NewInt(0), big.NewInt(0), big.NewInt(0)).Set(0, 0, 1)

	table := []struct {
		x, y, z  int64
		expected bool
	}{
		{0, 0, 0, true},
		{0, 0, 1, true},
		{0, 1, 0, true},
		{0, 1, 1, false},
		{1, 0, 0, true},
		{1, 0, 1, false},
		{1, 1, 0, true},
		{1, 1, 1, false},

		{64, 0, 0, true},
		{64, 0, 1, false},
		{64, 1, 0, true},
		{64, 1, 1, false},
		{65, 0, 0, true},
		{65, 0, 1, false},
		{65, 1, 0, true},
		{65, 1, 1, false},

		{0, 0, -64, true},
		{0, 0, -63, true},
		{0, 1, -64, true},
		{0, 1, -63, true},
		{1, 0, -64, true},
		{1, 0, -63, true},
		{1, 1, -64, true},
		{1, 1, -63, true},

		{0, 0, 64, false},
		{0, 0, 65, false},
		{0, 1, 64, false},
		{0, 1, 65, false},
		{1, 0, 64, false},
		{1, 0, 65, false},
		{1, 1, 64, false},
		{1, 1, 65, false},
	}

	for i, c := range table {
		if actual := w.Get(big.NewInt(c.x), big.NewInt(c.y), big.NewInt(c.z)); actual != c.expected {
			t.Errorf("case %d failed: %+v", i, c)
		}
	}
}
