package example

import "evm/util"

// Address is the address
type Address [20]byte

// Bytes is the implementation of interface
func (a *Address) Bytes() []byte {
	return a[:]
}

// BytesToAddress convert bytes to address
func BytesToAddress(bytes []byte) *Address {
	var a Address
	copy(a[:], util.FixBytesLength(bytes, 20))
	return &a
}

// HexToAddress convert hex string to address, string may begin with 0x, 0X or nothing
func HexToAddress(hex string) *Address {
	var a Address
	if bytes, err := util.HexToBytes(hex); err == nil {
		copy(a[:], util.FixBytesLength(bytes, 20))
	}
	return &a
}

// RandomAddress return a random address
func RandomAddress() *Address {
	var a Address
	for i := range a {
		a[i] = byte(util.RandNum(128))
	}
	return &a
}

// ZeroAddress return zero address
func ZeroAddress() *Address {
	var a Address
	return &a
}
