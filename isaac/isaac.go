/*
------------------------------------------------------------------------------
isaac.go: an implementation of Bob Jenkins' random number generator ISAAC based on 'readable.c'
* 18 Aug 2014 -- direct port of readable.c to Go
* 10 Sep 2014 -- updated to be more idiomatic Go
------------------------------------------------------------------------------
*/

package isaac

type isaac struct {
	/* external results */
	randrsl [256]uint32
	randcnt uint32

	/* internal state */
	mm         [256]uint32
	aa, bb, cc uint32
}

func (r *isaac) isaac() {
	var x, y uint32

	r.cc++       /* cc just gets incremented once per 256 results */
	r.bb += r.cc /* then combined with bb                         */

	for i := 0; i < 256; i++ {
		x = r.mm[i]

		switch i % 4 {
		case 0:
			r.aa ^= (r.aa << 13)
		case 1:
			r.aa ^= (r.aa >> 6)
		case 2:
			r.aa ^= (r.aa << 2)
		case 3:
			r.aa ^= (r.aa >> 16)
		}

		r.aa = r.mm[(i+128)%256] + r.aa
		y = r.mm[(x>>2)%256] + r.aa + r.bb
		r.bb = r.mm[(y>>10)%256] + x

		r.mm[i] = y
		r.randrsl[i] = r.bb

		/* Note that bits 2..9 are chosen from x but 10..17 are chosen
		   from y.  The only important thing here is that 2..9 and 10..17
		   don't overlap.  2..9 and 10..17 were then chosen for speed in
		   the optimized version (rand.c) */
		/* See http://burtleburtle.net/bob/rand/isaac.html
		   for further explanations and analysis. */

	}
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

/* if (flag==true), then use the contents of randrsl[] to initialize mm[]. */
func (r *isaac) randInit(flag bool) {
	var a, b, c, d, e, f, g, h uint32
	a, b, c, d, e, f, g, h = 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9

	for i := 0; i < 4; i++ {
		a, b, c, d, e, f, g, h = mix(a, b, c, d, e, f, g, h)
	}

	for i := 0; i < 256; i += 8 { /* fill mm[] with messy stuff */
		if flag { /* use all the information in the seed */
			a += r.randrsl[i]
			b += r.randrsl[i+1]
			c += r.randrsl[i+2]
			d += r.randrsl[i+3]
			e += r.randrsl[i+4]
			f += r.randrsl[i+5]
			g += r.randrsl[i+6]
			h += r.randrsl[i+7]
		}
		a, b, c, d, e, f, g, h = mix(a, b, c, d, e, f, g, h)
		r.mm[i] = a
		r.mm[i+1] = b
		r.mm[i+2] = c
		r.mm[i+3] = d
		r.mm[i+4] = e
		r.mm[i+5] = f
		r.mm[i+6] = g
		r.mm[i+7] = h
	}

	if flag { /* do a second pass to make all of the seed affect all of mm */
		for i := 0; i < 256; i += 8 {
			a += r.mm[i]
			b += r.mm[i+1]
			c += r.mm[i+2]
			d += r.mm[i+3]
			e += r.mm[i+4]
			f += r.mm[i+5]
			g += r.mm[i+6]
			h += r.mm[i+7]
			a, b, c, d, e, f, g, h = mix(a, b, c, d, e, f, g, h)
			r.mm[i] = a
			r.mm[i+1] = b
			r.mm[i+2] = c
			r.mm[i+3] = d
			r.mm[i+4] = e
			r.mm[i+5] = f
			r.mm[i+6] = g
			r.mm[i+7] = h
		}
	}

	r.isaac()     /* fill in the first set of results */
	r.randcnt = 0 /* reset the counter */
}

/* there is no official method for doing this, but just writing the key to the
 * state array is how the demo code does it */
func (r *isaac) Seed(key string) {
	for idx, c := range key {
		if idx == len(r.randrsl) {
			break
		}
		r.randrsl[idx] = uint32(c)
		r.randInit(true)
	}
}

/* retrieve the next number in the sequence */
func (r *isaac) Rand() uint32 {
	rnd := r.randrsl[r.randcnt]
	r.randcnt++
	if r.randcnt == uint32(len(r.randrsl)) {
		// reset for another 256
		r.isaac()
		r.randcnt = 0
	}
	return rnd
}
