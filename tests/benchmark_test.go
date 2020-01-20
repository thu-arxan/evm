package tests

import (
	"evm"
	"evm/db"
	"evm/example"
	"evm/util"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	benchmark = false
)

func TestBenchmarkAllSol(t *testing.T) {
	if !benchmark {
		return
	}
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var caller = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	var size = 10000
	// deploy balance sol
	bin, err := util.ReadBinFile(balanceBin)
	require.NoError(t, err)
	var gas uint64 = 1000000
	vm := evm.New(bc, memoryDB, &evm.Context{
		Input: bin,
		Value: 0,
		Gas:   &gas,
	})
	balanceCode, balanceAddress, err := vm.Create(caller)
	require.NoError(t, err)
	// deploy math sol
	bin, err = util.ReadBinFile(mathBin)
	require.NoError(t, err)
	gas = 1000000
	vm = evm.New(bc, memoryDB, &evm.Context{
		Input: bin,
		Value: 0,
		Gas:   &gas,
	})
	mathCode, mathAddress, err := vm.Create(caller)
	require.NoError(t, err)
	// deploy money sol
	bin, err = util.ReadBinFile(moneyBin)
	require.NoError(t, err)
	gas = 1000000
	vm = evm.New(bc, memoryDB, &evm.Context{
		Input: bin,
		Value: 0,
		Gas:   &gas,
	})
	moneyCode, moneyAddress, err := vm.Create(caller)
	require.NoError(t, err)
	// get of balance
	var payload = mustPack(balanceAbi, "get")
	var begin = time.Now()
	call(t, bc, memoryDB, caller, balanceAddress, balanceCode, payload, size)
	// set of balance
	payload = mustPack(balanceAbi, "set", "20")
	call(t, bc, memoryDB, caller, balanceAddress, balanceCode, payload, size)
	// info of balance
	payload = mustPack(balanceAbi, "info")
	call(t, bc, memoryDB, caller, balanceAddress, balanceCode, payload, size)
	// chaos of math
	payload = mustPack(mathAbi, "chaos")
	call(t, bc, memoryDB, caller, mathAddress, mathCode, payload, size)
	// get of money
	payload = mustPack(moneyAbi, "get")
	call(t, bc, memoryDB, caller, moneyAddress, moneyCode, payload, size)
	// add of money
	payload = mustPack(moneyAbi, "add")
	call(t, bc, memoryDB, caller, moneyAddress, moneyCode, payload, size)
	duration := time.Since(begin)
	fmt.Printf("%d times running cost %v\n", size, duration)
}

func call(t *testing.T, bc evm.Blockchain, db evm.DB, caller evm.Address, contract evm.Address, code, payload []byte, times int) {
	for i := 0; i < times; i++ {
		var gasQuota uint64 = 1000000
		var gas = gasQuota
		_, err := evm.New(bc, db, &evm.Context{
			Input: payload,
			Value: 0,
			Gas:   &gas,
		}).Call(caller, contract, code)
		require.NoError(t, err)
	}
}
