package core

// Lengths of hashes in bytes.
const (
	HashLength = 32
)

// Hash represents the 32 byte sha3-256 hash of arbitrary data.
type Hash [HashLength]byte
