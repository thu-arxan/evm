package core

import (
	"bytes"
	"encoding/binary"
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

// UnmarshalFromHexBytes unmarshal from hex bytes
func (w *Word256) UnmarshalFromHexBytes(hexBytes []byte) error {
	bs, err := hex.DecodeString(string(hexBytes))
	if err != nil {
		return err
	}
	copy(w[:], bs)
	return nil
}

// HexString return hex string of word256
func (w Word256) HexString() string {
	return hex.EncodeUpperToString(w[:])
}

// Copy copy the word256
// TODO: really???
func (w Word256) Copy() Word256 {
	return w
}

// Bytes return the bytes of word256
func (w Word256) Bytes() []byte {
	return w[:]
}

// Prefix return the prefix of word256
func (w Word256) Prefix(n int) []byte {
	return w[:n]
}

// Postfix return the postfix of word256
func (w Word256) Postfix(n int) []byte {
	return w[32-n:]
}

// Word160 get a Word160 embedded a Word256 and padded on the left (as it is for account addresses in EVM)
func (w Word256) Word160() (w160 Word160) {
	copy(w160[:], w[Word256Word160Delta:])
	return
}

// Address convert Word256 to Address
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

// Compare compare two word256
// TODO: return what?
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

// Size return size
func (w Word256) Size() int {
	return Word256Bytes
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

// Int64ToWord256 convert int64 to word256
func Int64ToWord256(i int64) Word256 {
	return BigIntToWord256(SignExtend(big.NewInt(i), Word256Bits))
}

// Int64FromWord256 convert word256 to int64
func Int64FromWord256(word Word256) int64 {
	return BigIntFromWord256(word).Int64()
}

// BigIntToWord256 convert big.Int to word256
func BigIntToWord256(x *big.Int) Word256 {
	return LeftPadWord256(U256(x).Bytes())
}

// BigIntFromWord256 convert word256 to big.Int
func BigIntFromWord256(word Word256) *big.Int {
	return S256(new(big.Int).SetBytes(word[:]))
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
