package tests

import (
	"evm"
	"evm/db"
	"evm/example"
	"evm/util"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

)

var (
	createBin = "sols/Create_sol_C.bin"
	createAbi = "sols/Create_sol_C.abi"
	CCode  []byte
	CAddress evm.Address
)

func TestCreateSol(t *testing.T) {
	binBytes, err := util.ReadBinFile(createBin)
	require.NoError(err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var origin = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	var exceptAddress = `cd234a471b72ba2f1ccf0a70fcaba648a5eecd8d`
	CCode, CAddress = deployContract(t, memoryDB, bc, origin, binBytes, exceptAddress, "", 321307)
	callInfo(t, memoryDB, bc, origin, mustParsePayload(createAbi, "createAndGetBalance", big.NewInt(44), big.NewInt(0)), 84357)
}

func callCreate(t *testing.T, db evm.DB, bc evm.Blockchain, caller evm.Address, payload []byte, gasCost uint64) {
	var gasQuota uint64 = 100000
	var gas = gasQuota
	output, err := evm.New(bc, db, &evm.Context{
		Input: payload,
		Value: 0,
		Gas: &gas,
	}).Call(caller, CAddress, CCode)
	require.NoError(t, err)
	if gasCost != 0 {
		require.EqualValues(t, gasCost, gasQuota - gas)
	}

	t.Log(output)
	t.Log(gasQuota - gas)
}