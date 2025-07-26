package isaac

import (
	"bytes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
)

// Stream is a cipher stream that implements the ISAAC algorithm when generating keys for each operation.
type Stream struct {
	*Rand
}

// XORKeyStream XORs each byte in the given slice with a byte from the cipher's key stream. Dst and src must overlap
// entirely or not at all.
func (s *Stream) XORKeyStream(dst, src []byte) {
	keyStream := new(bytes.Buffer)
	for len(src) > 0 {
		keyStream.Reset()

		// unpacking
		nextUint32 := s.Uint32()
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

// NewStream initializes a new cipher stream with the specified key.
func NewStream(key string) (cipher.Stream, error) {
	seed, err := TransformSeed(key)
	if err != nil {
		return nil, fmt.Errorf("failed to transform seed string %q: %w", key, err)
	}
	stream := &Stream{
		Rand: NewRand(seed...),
	}
	return stream, nil
}
