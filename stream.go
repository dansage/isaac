package isaac

import (
	"bytes"
	"crypto/cipher"
	"encoding/binary"
)

/* implementation based on http://golang.org/src/pkg/crypto/cipher/ctr.go */
func (r *ISAAC) XORKeyStream(dst, src []byte) {
	keyStream := new(bytes.Buffer)
	for len(src) > 0 {
		keyStream.Reset()

		// unpacking
		nextUint32 := r.Rand()
		binary.Write(keyStream, binary.BigEndian, &nextUint32)
		n := safeXORBytes(dst, src, keyStream.Bytes())

		dst = dst[n:]
		src = src[n:]
	}
}

func safeXORBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
}

func NewISAACStream(key string) cipher.Stream {
	stream := new(ISAAC)
	stream.Seed(key)
	return stream
}
