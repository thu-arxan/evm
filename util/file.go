package util

import (
	"io/ioutil"
)

// ReadBinFile read code from bin file
func ReadBinFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return HexToBytes(string(data))
}
