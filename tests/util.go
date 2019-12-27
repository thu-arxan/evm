package tests

import (
	"evm"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func deployContract(t *testing.T, db evm.DB, bc evm.Blockchain, caller evm.Address, bin []byte, exceptAddress, exceptCode string, gasCost uint64) ([]byte, evm.Address) {
	var originGas uint64 = 1000000
	if gasCost != 0 {
		originGas = gasCost
	}
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
