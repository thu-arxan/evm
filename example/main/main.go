//  Copyright 2020 The THU-Arxan Authors
//  This file is part of the evm library.
//
//  The evm library is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Lesser General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  The evm library is distributed in the hope that it will be useful,/
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
//  GNU Lesser General Public License for more details.
//
//  You should have received a copy of the GNU Lesser General Public License
//  along with the evm library. If not, see <http://www.gnu.org/licenses/>.
//

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
