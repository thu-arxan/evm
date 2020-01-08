package abi

import "fmt"

// Pack provide a easy way to pack, it is a simple wrapper of New & PackValues
func Pack(abiFile, funcName string, inputs ...string) ([]byte, error) {
	abi, err := New(abiFile)
	if err != nil {
		return nil, err
	}
	return abi.PackValues(funcName, inputs...)
}

// Unpack provide a easy way to unpack, it is a simple wrapper of New & UnpackValues
func Unpack(abiFile, funcName string, data []byte) (values []string, err error) {
	defer func() {
		if e := recover(); e != nil {
			values = nil
			err = fmt.Errorf("unpack failed because %v", e)
		}
	}()
	abi, err := New(abiFile)
	if err != nil {
		return nil, err
	}
	values, err = abi.UnpackValues(funcName, data)
	return
}
