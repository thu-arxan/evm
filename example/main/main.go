package main

import (
	"evm"
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
	var gas uint64
	gas = 10000000
	vm := evm.New(bc, db.NewMemory(bc.NewAccount), &evm.Context{
		Input: code,
		Value: 0,
		Gas:   &gas,
	})

	var origin = example.RandomAddress()
	code, _, err = vm.Create(origin)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", code)
}
