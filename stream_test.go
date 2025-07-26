package isaac

import (
	"bytes"
	"testing"
)

func TestStream(t *testing.T) {
	var key = "This is <i>not</i> the right mytext."
	var plaintext = []byte("Hello, world")
	var ciphertext = make([]byte, len(plaintext))
	var decrypted = make([]byte, len(plaintext))

	enc := NewISAACStream(key)
	enc.XORKeyStream(ciphertext, plaintext)

	dec := NewISAACStream(key)
	dec.XORKeyStream(decrypted, ciphertext)

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Plaintext not equal to decrypted ciphertext")
	}

}
