package main

import (
	"compress/gzip"
	"encoding/gob"
	"flag"
	"log"
	"math/big"
	"math/rand"
	"os"

	"github.com/BenLubar/terrain/world"
)

type World struct {
	*world.World
	Diggers []*Digger
}

type Digger struct {
	X, Y, Z big.Int
}

func main() {
	var (
		in         = flag.String("i", "", "read this file instead of making a new world")
		out        = flag.String("o", "world.gz", "save the world to this location at the end of execution")
		iterations = flag.Uint64("n", 10000, "number of simulation iterations")
		seed       = flag.Int64("s", 0, "random seed")
	)

	flag.Parse()
	if len(flag.Args()) != 0 {
		flag.Usage()
		os.Exit(1)
	}

	var w World
	if *in == "" {
		w.World = world.New()
	} else {
		f, err := os.Open(*in)
		if err != nil {
			log.Fatalf("error opening input file %q: %v", *in, err)
		}
		g, err := gzip.NewReader(f)
		if err != nil {
			f.Close()
			log.Fatalf("error decompressing input file %q: %v", *in, err)
		}
		err = gob.NewDecoder(g).Decode(&w)
		g.Close()
		f.Close()
		if err != nil {
			log.Fatalf("error decoding input file %q: %v", *in, err)
		}
	}

	one := big.NewInt(1)
	var scratch big.Int
	r := rand.New(rand.NewSource(*seed))
	for i := uint64(0); i < *iterations; i++ {
		switch r.Intn(3) {
		case 0, 1:
			if len(w.Diggers) != 0 {
				d := w.Diggers[r.Intn(len(w.Diggers))]
				if !w.Get(&d.X, &d.Y, scratch.Sub(&d.Z, one)) {
					// Digger fell down a tunnel
					d.Z.Set(&scratch)
					continue
				}
				switch r.Intn(17) {
				case 0:
					d.X.Add(&d.X, one)
				case 1:
					d.X.Add(&d.X, one)
					d.Y.Add(&d.Y, one)
				case 2:
					d.Y.Add(&d.Y, one)
				case 3:
					d.X.Sub(&d.X, one)
					d.Y.Add(&d.Y, one)
				case 4:
					d.X.Sub(&d.X, one)
				case 5:
					d.X.Sub(&d.X, one)
					d.Y.Sub(&d.Y, one)
				case 6:
					d.Y.Sub(&d.Y, one)
				case 7:
					d.X.Add(&d.X, one)
					d.Y.Sub(&d.Y, one)
				case 8:
					d.X.Add(&d.X, one)
					d.Z.Sub(&d.Z, one)
				case 9:
					d.X.Add(&d.X, one)
					d.Y.Add(&d.Y, one)
					d.Z.Sub(&d.Z, one)
				case 10:
					d.Y.Add(&d.Y, one)
					d.Z.Sub(&d.Z, one)
				case 11:
					d.X.Sub(&d.X, one)
					d.Y.Add(&d.Y, one)
					d.Z.Sub(&d.Z, one)
				case 12:
					d.X.Sub(&d.X, one)
					d.Z.Sub(&d.Z, one)
				case 13:
					d.X.Sub(&d.X, one)
					d.Y.Sub(&d.Y, one)
					d.Z.Sub(&d.Z, one)
				case 14:
					d.Y.Sub(&d.Y, one)
					d.Z.Sub(&d.Z, one)
				case 15:
					d.X.Add(&d.X, one)
					d.Y.Sub(&d.Y, one)
					d.Z.Sub(&d.Z, one)
				case 16:
					d.Z.Sub(&d.Z, one)
				}
				w.Unset(&d.X, &d.Y, &d.Z)
				continue
			}
			fallthrough
		case 2:
			w.Diggers = append(w.Diggers, &Digger{})
		}
	}

	f, err := os.Create(*out)
	if err != nil {
		log.Fatalf("error creating output file %q: %v", *out, err)
	}
	defer f.Close()
	g, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		log.Fatalf("error compressing output file %q: %v", *out, err)
	}
	defer g.Close()
	err = gob.NewEncoder(g).Encode(&w)
	if err != nil {
		log.Fatalf("error encoding output file %q: %v", *out, err)
	}
}
