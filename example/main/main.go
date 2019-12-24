package main

import (
	"evm"
	"evm/example"
	eutil "evm/example/util"
	"fmt"
)

func main() {
	code, err := eutil.ReadBinFile("../sols/output/Balance_sol_Balance.bin")
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
		Value:  0,
		Gas:    &gas,
	}, code)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", code)
}
