package core

import (
	"math/big"
)

var Big1 = big.NewInt(1)

// Big256 is the big.Int of 256
var Big256 = big.NewInt(256)

// U256 converts a possibly negative big int x into a positive big int encoding a twos complement representation of x
// truncated to 32 bytes
func U256(x *big.Int) *big.Int {
	return ToTwosComplement(x, Word256Bits)
}

// S256 interprets a positive big.Int as a 256-bit two's complement signed integer
func S256(x *big.Int) *big.Int {
	return FromTwosComplement(x, Word256Bits)
}

// ToTwosComplement convert a possibly negative big.Int x to a positive big.Int encoded in two's complement
func ToTwosComplement(x *big.Int, n uint) *big.Int {
	// And treats negative arguments a if they were twos complement encoded so we end up with a positive number here
	// with the twos complement bit pattern
	return new(big.Int).And(x, andMask(n))
}

// FromTwosComplement interprets a positive big.Int as a n-bit two's complement signed integer
func FromTwosComplement(x *big.Int, n uint) *big.Int {
	signBit := int(n) - 1
	if x.Bit(signBit) == 0 {
		// Sign bit not set => value (v) is positive
		// x = |v| = v
		return x
	}
	// Sign bit set => value (v) is negative
	// x = 2^n - |v|
	b := new(big.Int).Lsh(Big1, n)
	// v = -|v| = x - 2^n
	return new(big.Int).Sub(x, b)
}

// SignExtend treats the positive big int x as if it contains an embedded n bit signed integer in its least significant
// bits and extends that sign
func SignExtend(x *big.Int, n uint) *big.Int {
	signBit := n - 1
	// single bit set at sign bit position
	mask := new(big.Int).Lsh(Big1, signBit)
	// all bits below sign bit set to 1 all above (including sign bit) set to 0
	mask.Sub(mask, Big1)
	if x.Bit(int(signBit)) == 1 {
		// Number represented is negative - set all bits above sign bit (including sign bit)
		return x.Or(x, mask.Not(mask))
	}
	// Number represented is positive - clear all bits above sign bit (including sign bit)
	return x.And(x, mask)
}

func andMask(n uint) *big.Int {
	x := new(big.Int)
	return x.Sub(x.Lsh(Big1, n), Big1)
}
