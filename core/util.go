package core

import "evm/util"

func randBytes(length int) []byte {
	var bytes = make([]byte, length)
	for i := range bytes {
		bytes[i] = byte(util.RandNum(10))
	}
	return bytes
}
