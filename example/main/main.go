package main

import (
	"evm"
	"evm/example"
	"evm/util"
	"fmt"
	"io/ioutil"
)

func main() {
	code, err := readBinCode("../sols/output/Balance_sol_Balance.bin")
	if err != nil {
		panic(err)
	}
	vm := evm.New(example.NewBlockchain(), example.NewMemoryDB())
	vm.Create(evm.Context{}, code)
}

func readBinCode(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))
	return util.HexToBytes(string(data))
}
