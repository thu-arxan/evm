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
	var gas uint64
	gas = 10000000
	var origin = example.RandomAddress()
	code, err = vm.Call(evm.Context{
		Origin: origin,
		Caller: origin,
		Callee: example.ZeroAddress(),
		Value:  10,
		Gas:    &gas,
	}, code)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", code)
}

func readBinCode(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(data))
	return util.HexToBytes(string(data))
}
