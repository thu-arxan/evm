package util

import (
	"evm/util"
	"io/ioutil"
)

// ReadBinFile read code from bin file
func ReadBinFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return util.HexToBytes(string(data))
}
