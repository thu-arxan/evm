package core

import "fmt"

// Lengths of hashes in bytes.
const (
	HashLength = 32
)

// Hash represents the 32 byte sha3-256 hash of arbitrary data.
type Hash [HashLength]byte

// BytesToHash convert bytes to hash
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// SetBytes Sets the hash to the value of b. If b is larger than len(h), 'b' will be cropped (from the left).
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// Bytes return bytes of hash
func (h Hash) Bytes() []byte {
	return h[:]
}

// Hex return hex of hash
func (h Hash) Hex() string {
	return fmt.Sprintf("%x", h.Bytes())
}
