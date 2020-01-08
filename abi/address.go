package abi

// This file provide a function which allow user define the address unpack they need.

var (
	addressLength   = 20
	addressToString func([]byte) string
)

// SetAddressParser set address reflect type and its length and toString function
func SetAddressParser(length int, toString func([]byte) string) {
	addressLength = length
	addressToString = toString
}

// ResetAddressParser reset reflect type and its length and toString function
func ResetAddressParser() {
	addressLength = 20
	addressToString = nil
}
