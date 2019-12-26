package tests

import (
	"evm/util"
	"fmt"
)

// AddressLength define the length of Address
const AddressLength = 12

// Address is the address
type Address struct {
	value [AddressLength]byte
}

// HexToAddress convert hex string to address
func HexToAddress(s string) (*Address, error) {
	bytes, err := util.HexToBytes(s)
	if err != nil {
		return nil, err
	}
	bytes = util.FixBytesLength(bytes, AddressLength)
	var a Address
	for i := range a.value {
		a.value[i] = bytes[i]
	}
	return &a, nil
}

// BytesToAddress convert bytes to address
func BytesToAddress(bytes []byte) *Address {
	bytes = util.FixBytesLength(bytes, AddressLength)
	var a Address
	for i := range a.value {
		a.value[i] = bytes[i]
	}
	return &a
}

// Bytes return bytes of Address
func (a *Address) Bytes() []byte {
	var bytes = make([]byte, AddressLength)
	for i := range bytes {
		bytes[i] = a.value[i]
	}
	return bytes
}

// Length return the length of address
func (a *Address) Length() int {
	return AddressLength
}

// RandomAddress random generate an address
func RandomAddress() *Address {
	var a Address
	for i := range a.value {
		a.value[i] = byte(util.RandNum(128))
	}
	return &a
}

// String return hex of address
func (a *Address) String() string {
	return fmt.Sprintf("%X", a.value)
}
