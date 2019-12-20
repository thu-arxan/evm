package util

import (
	"encoding/binary"
	"fmt"
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
