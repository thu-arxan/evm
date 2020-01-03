package abi

import (
	"reflect"
)

// This file provide a function which allow user define the address unpack they need.

var (
	addressLength   = 20
	addressToString func([]byte) string
)

// SetAddressParser set address reflect type and its length
func SetAddressParser(t reflect.Type, length int) {
	addressT = t
	addressLength = length
}
