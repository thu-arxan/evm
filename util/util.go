package util

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"strings"
	"time"
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

// HexToBytes will remove 0x or 0X of begin ,and then call hex.DecodeString
func HexToBytes(s string) ([]byte, error) {
	if strings.HasPrefix(s, "0x") {
		s = strings.Replace(s, "0x", "", 1)
	} else if strings.HasPrefix(s, "0X") {
		s = strings.Replace(s, "0X", "", 1)
	}
	return hex.DecodeString(s)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
// Note: The function is a wrapper of hex.DecodeString, but it ignore the error.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

// RandNum return int in [0, max)
func RandNum(max int) int {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max)
	return randNum
}

// BytesCombine combines some bytes
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// FixBytesLength fix bytes to bytes which length is length
func FixBytesLength(bytes []byte, length int) []byte {
	var result = make([]byte, length)
	if len(bytes) > length {
		bytes = bytes[len(bytes)-length:]
	} else if len(bytes) < length {
		bytes = BytesCombine(make([]byte, length-len(bytes)), bytes)
	}
	copy(result[:], bytes)
	return result
}

// Hex is the wrapper of fmt.Sprintf("%x", data)
func Hex(data []byte) string {
	return fmt.Sprintf("%x", data)
}

// Log256 call and down round
func Log256(x *big.Int) int {
	if x.Sign() <= 0 {
		return 0
	}
	return (len(x.Text(2)) - 1) / 8
}
