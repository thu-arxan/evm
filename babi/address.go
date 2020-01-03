package abi

import "errors"

// This file provide a function which allow user define the address unpack they need.

var (
	addressLength   = 20
	addressToString func([]byte) string
)

// SetAddressParser set address length and toString function
// If length > 20, something may be lost while running
func SetAddressParser(length int, toString func([]byte) string) error {
	if length <= 0 || length > 20 {
		return errors.New("length should belong to (0,20]")
	}
	addressLength = length
	addressToString = toString
	return nil
}

// ResetAddressParser reset to default
func ResetAddressParser() {
	addressLength = 20
	addressToString = nil
}
