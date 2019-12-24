package core

import (
	"evm/util"
	"strings"
)

// Address is Word160
type Address Word160

// Here defines some length
const (
	AddressLength    = Word160Length
	AddressHexLength = 2 * AddressLength
)

// ZeroAddress is zero
var ZeroAddress = Address{}

// Word256 return Word256 of address
func (address Address) Word256() Word256 {
	return Word160(address).Word256()
}

// Bytes return bytes of address
func (address Address) Bytes() []byte {
	return Word160(address).Bytes()
}

// AddressFromBytes returns an addres. It will cut left if len(bs) > 20, else add zeros at left if len(bs) < 20
func AddressFromBytes(bs []byte) (address Address) {
	if len(bs) > Word160Length {
		bs = bs[len(bs)-Word160Length:]
	}
	copy(address[:], bs)
	return
}

// HexToAddress convert hex string to Address, hex string could begin with 0x or 0X or nothing
func HexToAddress(s string) (address Address, err error) {
	if strings.HasPrefix(s, "0x") {
		s = strings.Replace(s, "0x", "", 1)
	} else if strings.HasPrefix(s, "0X") {
		s = strings.Replace(s, "0X", "", 1)
	}
	bytes, err := util.HexToBytes(s)
	if err != nil {
		return ZeroAddress, err
	}
	return AddressFromBytes(bytes), nil
}
