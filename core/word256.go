package core

import (
	"bytes"
	"encoding/binary"
	"evm/util"
	"math/big"

	hex "github.com/tmthrgd/go-hex"
)

// Here defines some word256
var (
	Zero256 = Word256{}
	One256  = LeftPadWord256([]byte{1})
)

// Here defines some length
const (
	Word256Bytes = 32
	Word256Bits  = Word256Bytes * 8
)

// Here defines some variables
var (
	BigWord256Bytes = big.NewInt(Word256Bytes)
	trimCutSet      = string([]byte{0})
)

// Word256 is bytes which length is Word256Bytes
type Word256 [Word256Bytes]byte

// HexString return hex string of word256
func (w Word256) HexString() string {
	return hex.EncodeUpperToString(w[:])
}

// Copy copy the word256
func (w Word256) Copy() Word256 {
	return w
}

// Bytes return the bytes of word256
func (w Word256) Bytes() []byte {
	return w[:]
}

// Word160 get a Word160 embedded a Word256
// It will remove left zeros until length == 20
func (w Word256) Word160() (w160 Word160) {
	copy(w160[:], w[Word256Word160Delta:])
	return
}

// Address convert Word256 to Address
// It is a wrapper of Word160(which is same as EVM)
func (w Word256) Address() Address {
	return Address(w.Word160())
}

// IsZero return if word256 is zero
func (w Word256) IsZero() bool {
	accum := byte(0)
	for _, byt := range w {
		accum |= byt
	}
	return accum == 0
}

// Compare compare two word256, it will return 0 if a == b
func (w Word256) Compare(other Word256) int {
	return bytes.Compare(w[:], other[:])
}

// UnpadLeft trim left zeros
func (w Word256) UnpadLeft() []byte {
	return bytes.TrimLeft(w[:], trimCutSet)
}

// UnpadRight trim right zeros
func (w Word256) UnpadRight() []byte {
	return bytes.TrimRight(w[:], trimCutSet)
}

// Uint64ToWord256 convert uint64 to word256
func Uint64ToWord256(i uint64) (word Word256) {
	binary.BigEndian.PutUint64(word[24:], i)
	return
}

// Uint64FromWord256 convert word256 to uint64
func Uint64FromWord256(word Word256) uint64 {
	return binary.BigEndian.Uint64(word.Postfix(8))
}

// RightPadWord256 keep the right pad of word256
func RightPadWord256(bz []byte) (word Word256) {
	copy(word[:], bz)
	return
}

// LeftPadWord256 keep the left pad of word256
func LeftPadWord256(bz []byte) (word Word256) {
	copy(word[32-len(bz):], bz)
	return
}

// Postfix return the postfix of word256
func (w Word256) Postfix(n int) []byte {
	return w[32-n:]
}

// BytesToWord256 convert bytes to Word256.
// It will add left zero if len(data) < 32, and remove left bytes if len(data) > 32
func BytesToWord256(bs []byte) (word Word256) {
	if len(bs) > Word256Bytes {
		bs = bs[len(bs)-Word256Bytes:]
	} else if len(bs) < Word256Bytes {
		bs = util.BytesCombine(make([]byte, Word256Bytes-len(bs)), bs)
	}
	copy(word[:], bs)
	return
}
