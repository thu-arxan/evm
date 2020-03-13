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

package tests

import (
	"github.com/thu-arxan/evm"
	abi "github.com/thu-arxan/evm/abi"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func deployContract(t *testing.T, db evm.DB, bc evm.Blockchain, caller evm.Address, bin []byte, exceptAddress, exceptCode string, gasCost uint64) ([]byte, evm.Address) {
	var originGas uint64 = 1000000
	var gas = originGas
	vm := evm.New(bc, db, &evm.Context{
		Input: bin,
		Value: 0,
		Gas:   &gas,
	})
	code, address, err := vm.Create(caller)
	require.NoError(t, err)
	if gasCost != 0 {
		require.EqualValues(t, gasCost, originGas-gas, fmt.Sprintf("except %d other than %d", gasCost, originGas-gas))
	}
	if exceptCode != "" {
		require.Equal(t, exceptCode, fmt.Sprintf("%x", code))
	}
	if exceptAddress != "" {
		require.Equal(t, exceptAddress, fmt.Sprintf("%x", address.Bytes()))
	}
	account := db.GetAccount(address)
	require.Equal(t, code, account.GetCode())
	return code, address
}

func deployContractWithValue(t *testing.T, db evm.DB, bc evm.Blockchain, caller evm.Address, bin []byte, value, gasCost uint64) ([]byte, evm.Address) {
	var originGas uint64 = 1000000
	var gas = originGas
	vm := evm.New(bc, db, &evm.Context{
		Input: bin,
		Value: value,
		Gas:   &gas,
	})
	code, address, err := vm.Create(caller)
	require.NoError(t, err)
	if gasCost != 0 {
		require.EqualValues(t, gasCost, originGas-gas, fmt.Sprintf("except %d other than %d", gasCost, originGas-gas))
	}
	account := db.GetAccount(address)
	require.Equal(t, code, account.GetCode())
	require.Equal(t, uint64(1), account.GetNonce())
	return code, address
}

func callWithPayload(t *testing.T, db evm.DB, bc evm.Blockchain, caller, contract evm.Address, payload []byte, gasCost, refund uint64) []byte {
	var gasQuota uint64 = 100000
	var gas = gasQuota
	vm := evm.New(bc, db, &evm.Context{
		Input: payload,
		Value: 0,
		Gas:   &gas,
	})
	code := db.GetAccount(contract).GetCode()
	result, err := vm.Call(caller, contract, code)
	require.NoError(t, err)
	if gasCost != 0 {
		require.EqualValues(t, gasCost, gasQuota-gas, fmt.Sprintf("Except gas cost %d other than %d", gasCost, gasQuota-gas))
	}
	require.EqualValues(t, refund, vm.GetRefund(), fmt.Sprintf("Except refund %d other than %d", refund, vm.GetRefund()))
	return result
}

func callWithPayloadAndValue(t *testing.T, db evm.DB, bc evm.Blockchain, caller, contract evm.Address, payload []byte, value uint64, gasCost, refund uint64) []byte {
	var gasQuota uint64 = 100000
	var gas = gasQuota
	vm := evm.New(bc, db, &evm.Context{
		Input: payload,
		Value: value,
		Gas:   &gas,
	})
	code := db.GetAccount(contract).GetCode()
	result, err := vm.Call(caller, contract, code)
	require.NoError(t, err)
	if gasCost != 0 {
		require.EqualValues(t, gasCost, gasQuota-gas, fmt.Sprintf("Except gas cost %d other than %d", gasCost, gasQuota-gas))
	}
	require.EqualValues(t, refund, vm.GetRefund(), fmt.Sprintf("Except refund %d other than %d", refund, vm.GetRefund()))
	return result
}

func parsePayload(abiFile string, funcName string, args ...interface{}) ([]byte, error) {
	abi, err := abi.New(abiFile)
	if err != nil {
		return nil, err
	}
	return abi.Pack(funcName, args...)
}

func mustPack(abiFile, funcName string, inputs ...string) []byte {
	values, err := abi.Pack(abiFile, funcName, inputs...)
	if err != nil {
		panic(err)
	}
	fmt.Printf("pack payload is %x\n", values)
	return values
}

func mustUnpack(abiFile string, funcName string, data []byte) []string {
	values, err := abi.Unpack(abiFile, funcName, data)
	if err != nil {
		panic(err)
	}
	return values
}
