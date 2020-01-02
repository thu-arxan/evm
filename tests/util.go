package tests

import (
	"evm"
	"evm/abi"
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
	require.Equal(t, uint64(1), account.GetNonce())
	return code, address
}

func call(t *testing.T, db evm.DB, bc evm.Blockchain, caller, contract evm.Address, abiFile, funcName string, inputs, excepts []string, gasCost, refund uint64) {
	payload, err := abi.GetPayloadBytes(abiFile, funcName, inputs)
	require.NoError(t, err)
	var gasQuota uint64 = 100000
	var gas = gasQuota
	vm := evm.New(bc, db, &evm.Context{
		Input: payload,
		Value: 0,
		Gas:   &gas,
	})
	code := db.GetAccount(contract).GetCode()
	output, err := vm.Call(caller, contract, code)
	require.NoError(t, err)
	variables, err := abi.Unpacker(abiFile, funcName, output)
	require.NoError(t, err)
	require.Len(t, variables, len(excepts))
	for i := range excepts {
		require.Equal(t, excepts[i], variables[i].Value)
	}
	if gasCost != 0 {
		require.EqualValues(t, gasCost, gasQuota-gas, fmt.Sprintf("Except gas cost %d other than %d", gasCost, gasQuota-gas))
	}
	require.EqualValues(t, refund, vm.GetRefund(), fmt.Sprintf("Except refund %d other than %d", refund, vm.GetRefund()))
}
