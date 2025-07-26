package isaac

import (
	"bytes"
	"encoding/binary"
)

// Size is the number of entries in the internal memory and results slices.
const Size uint32 = 256

// Rand is a pseudorandom number generator that implements the ISAAC algorithm.
type Rand struct {
	/* external results */
	// results is the slice of generated pseudorandom numbers to be returned to callers.
	results []uint32

	// index is the index of the next number to return from the results slice.
	index uint32

	// memory is the internal state of the generator.
	memory []uint32

	// accumulator is used to maintain the state of the generator throughout each operation.
	accumulator uint32

	// lastResult is the last pseudorandom number generated.
	lastResult uint32

	// counter tracks the number of times a set of pseudorandom numbers has been generated. This is incremented every Size
	// numbers.
	counter uint32
}

// NewRand initializes a new instance of the pseudorandom number generator, injecting the seed values if specified
func NewRand(seed ...uint32) *Rand {
	i := new(Rand)

	// initialize the slices in the correct size
	i.memory = make([]uint32, Size)
	i.results = make([]uint32, Size)

	// if provided, inject the specified seed into the results
	if len(seed) > 0 {
		copy(i.results, seed)
	}

	// initialize the internal memory of the pseudorandom number generator
	i.randInit(true)
	return i
}

func (r *Rand) isaac() {
	// increment the counter and add it to the last result
	r.counter = r.counter + 1
	r.lastResult = r.lastResult + r.counter

	for i := uint32(0); i < Size; i++ {
		x := r.memory[i]
		switch i % 4 {
		case 0:
			r.accumulator = r.accumulator ^ (r.accumulator << 13)
		case 1:
			r.accumulator = r.accumulator ^ (r.accumulator >> 6)
		case 2:
			r.accumulator = r.accumulator ^ (r.accumulator << 2)
		case 3:
			r.accumulator = r.accumulator ^ (r.accumulator >> 16)
		}
		r.accumulator = r.memory[(i+128)%Size] + r.accumulator
		y := r.memory[(x>>2)%Size] + r.accumulator + r.lastResult
		r.memory[i] = y
		r.lastResult = r.memory[(y>>10)%Size] + x
		r.results[i] = r.lastResult
	}

	/* Note that bits 2..9 are chosen from x but 10..17 are chosen
	   from y.  The only important thing here is that 2..9 and 10..17
	   don't overlap.  2..9 and 10..17 were then chosen for speed in
	   the optimized version (rand.c) */
	/* See http://burtleburtle.net/bob/rand/isaac.html
	   for further explanations and analysis. */
}

func mix(a, b, c, d, e, f, g, h uint32) (uint32, uint32, uint32, uint32, uint32, uint32, uint32, uint32) {
	a ^= b << 11
	d += a
	b += c
	b ^= c >> 2
	e += b
	c += d
	c ^= d << 8
	f += c
	d += e
	d ^= e >> 16
	g += d
	e += f
	e ^= f << 10
	h += e
	f += g
	f ^= g >> 4
	a += f
	g += h
	g ^= h << 8
	b += g
	h += a
	h ^= a >> 9
	c += h
	a += b
	return a, b, c, d, e, f, g, h
}

/* if (flag==true), then use the contents of results[] to initialize memory[]. */
func (r *Rand) randInit(flag bool) {
	var a, b, c, d, e, f, g, h uint32
	a, b, c, d, e, f, g, h = 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9

	for i := 0; i < 4; i++ {
		a, b, c, d, e, f, g, h = mix(a, b, c, d, e, f, g, h)
	}

	for i := uint32(0); i < Size; i += 8 { /* fill memory[] with messy stuff */
		if flag { /* use all the information in the seed */
			a += r.results[i]
			b += r.results[i+1]
			c += r.results[i+2]
			d += r.results[i+3]
			e += r.results[i+4]
			f += r.results[i+5]
			g += r.results[i+6]
			h += r.results[i+7]
		}
		a, b, c, d, e, f, g, h = mix(a, b, c, d, e, f, g, h)
		r.memory[i] = a
		r.memory[i+1] = b
		r.memory[i+2] = c
		r.memory[i+3] = d
		r.memory[i+4] = e
		r.memory[i+5] = f
		r.memory[i+6] = g
		r.memory[i+7] = h
	}

	if flag { /* do a second pass to make all of the seed affect all of memory */
		for i := uint32(0); i < Size; i += 8 {
			a += r.memory[i]
			b += r.memory[i+1]
			c += r.memory[i+2]
			d += r.memory[i+3]
			e += r.memory[i+4]
			f += r.memory[i+5]
			g += r.memory[i+6]
			h += r.memory[i+7]
			a, b, c, d, e, f, g, h = mix(a, b, c, d, e, f, g, h)
			r.memory[i] = a
			r.memory[i+1] = b
			r.memory[i+2] = c
			r.memory[i+3] = d
			r.memory[i+4] = e
			r.memory[i+5] = f
			r.memory[i+6] = g
			r.memory[i+7] = h
		}
	}

	// generate the first set of results and reset the index
	r.isaac()
	r.index = Size
}

// TransformSeed creates a properly padded seed value from the specified seed string.
func TransformSeed(s string) ([]uint32, error) {
	// convert the specified seed into a byte slice
	sb := bytes.NewBuffer([]byte(s))

	// ensure the byte slice length is a multiple of 4 to represent a slice of uint32 values
	if sb.Len()%4 != 0 {
		var padding = 4 - (sb.Len() % 4)
		for i := 0; i < padding; i++ {
			// add a null byte to the buffer as padding
			sb.WriteByte(0x0)
		}
	}

	// loop through each possible uint32 value in the seed
	var value uint32
	var seed []uint32
	var values = sb.Len() / 4
	for i := 0; i < values; i++ {
		// convert the value into a uint32 value
		if err := binary.Read(sb, binary.LittleEndian, &value); err != nil {
			return nil, err
		}

		// add the value to the seed slice
		seed = append(seed, value)
	}
	return seed, nil
}

// Uint32 returns a pseudorandom 32-bit value as a uint32.
func (r *Rand) Uint32() (n uint32) {
	// decrement the index and prepare the result
	r.index--
	n = r.results[r.index]

	// check if the last value has been reached
	if r.index == 0 {
		// generate new values and reset the index
		r.isaac()
		r.index = Size
	}
	return n
}
