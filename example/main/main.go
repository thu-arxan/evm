package main

import (
	"evm"
	"evm/abi"
	"evm/db"
	"evm/example"
	"evm/util"
	"fmt"
)

func main() {
	code, err := util.ReadBinFile("../sols/output/Balance_sol_Balance.bin")
	if err != nil {
		panic(err)
	}
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var gas uint64
	gas = 10000000
	vm := evm.New(bc, memoryDB, &evm.Context{
		Input: code,
		Value: 0,
		Gas:   &gas,
	})

	var caller = example.RandomAddress()
	code, callee, err := vm.Create(caller)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", code)

	gas = 1000000
	payload, err := abi.Pack("../sols/output/Balance_sol_Balance.abi", "get")
	if err != nil {
		panic(err)
	}
	output, err := evm.New(bc, memoryDB, &evm.Context{
		Input: payload,
		Value: 0,
		Gas:   &gas,
	}).Call(caller, callee, code)
	fmt.Printf("%x\n", output)
}
