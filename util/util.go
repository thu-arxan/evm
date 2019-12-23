package util

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"
)

// SubSlice returns a subslice from offset of length length and a bool
// (true iff slice was possible). If the subslice
// extends past the end of data it returns A COPY of the segment at the end of
// data padded with zeroes on the right. If offset == len(data) it returns all
// zeroes. if offset > len(data) it returns a false
func SubSlice(data []byte, offset, length uint64) ([]byte, error) {
	size := uint64(len(data))
	if size < offset || offset < 0 || length < 0 {
		return nil, fmt.Errorf("data size is %d while offset and length are %d and %d", size, offset, length)
		// return nil, errors.Errorf(errors.Codes.InputOutOfBounds,
		// 	"subslice could not slice data of size %d at offset %d for length %d", size, offset, length)
	}
	if size < offset+length {
		// Extract slice from offset to end padding to requested length
		ret := make([]byte, length)
		copy(ret, data[offset:])
		return ret, nil
	}
	return data[offset : offset+length], nil
}

// Uint64ToBytes turn int64 to []byte
func Uint64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

// RightPadBytes ...
func RightPadBytes(slice []byte, l int) []byte {
	if l < len(slice) {
		return slice
	}
	padded := make([]byte, l)
	copy(padded[0:len(slice)], slice)
	return padded
}

// LeftPadBytes ...
func LeftPadBytes(slice []byte, l int) []byte {
	if l < len(slice) {
		return slice
	}
	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)
	return padded
}

// Contain return if the target which is a map or slice contains the obj
func Contain(target interface{}, obj interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

// HexToBytes is the wrapper of hex.DecodeString
func HexToBytes(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
