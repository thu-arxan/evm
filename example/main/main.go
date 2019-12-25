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
	vm := evm.New(bc, db.NewMemory(bc.NewAccount))
	var gas uint64
	gas = 10000000
	var origin = example.RandomAddress()
	code, err = vm.Call(evm.Context{
		Origin: origin,
		Caller: origin,
		Callee: example.ZeroAddress(),
		Value:  0,
		Gas:    &gas,
	}, code)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", code)
}
