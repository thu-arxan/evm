package crypto

import (
	"crypto/sha256"

	"golang.org/x/crypto/sha3"
)

// Keccak256 use sha3 to hash data
func Keccak256(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}

// SHA256 use sha256 to hash data
func SHA256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
