package abi

// Pack provide a easy way to pack, it is a simple wrapper of New & PackValues
func Pack(abiFile, funcName string, inputs ...string) ([]byte, error) {
	abi, err := New(abiFile)
	if err != nil {
		return nil, err
	}
	return abi.PackValues(funcName, inputs...)
}
