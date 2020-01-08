package abi

// This file provide a function which allow user define the address unpack they need.

var (
	addressLength   = 20
	addressToString func([]byte) string
	stringToAddress func(string) ([]byte, error)
)

// SetAddressParser set address reflect type and its length and toString function
func SetAddressParser(length int, toString func([]byte) string, toAddress func(string) ([]byte, error)) {
	addressLength = length
	addressToString = toString
	stringToAddress = toAddress
}

// ResetAddressParser reset reflect type and its length and toString function
func ResetAddressParser() {
	addressLength = 20
	addressToString = nil
	stringToAddress = nil
}
