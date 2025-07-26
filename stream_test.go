package isaac

import (
	"bytes"
	"testing"
)

func TestStream(t *testing.T) {
	var key = "This is <i>not</i> the right mytext."
	var message = []byte("Hello, world")

	// create a stream to use for encrypting the specified string
	enc, err := NewStream(key)
	if err != nil {
		t.Fatalf("failed to create new ISAAC stream: %v", err)
	}

	// encrypt the message using the stream
	encrypted := make([]byte, len(message))
	enc.XORKeyStream(encrypted, message)

	// create a stream to use for decrypting the specified string
	dec, err := NewStream(key)
	if err != nil {
		t.Fatalf("failed to create new ISAAC stream: %v", err)
	}

	// encrypt the message using the stream
	decrypted := make([]byte, len(message))
	dec.XORKeyStream(decrypted, encrypted)

	// ensure the messages are identical
	if !bytes.Equal(decrypted, message) {
		t.Fatalf("decrypted message doesn't match original message, expected %v got %v", message, decrypted)
	}
}
