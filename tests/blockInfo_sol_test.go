package tests

import (
	"evm"
	"evm/db"
	"evm/example"
	"evm/util"

	"testing"
	"github.com/stretchr/testify/require"
)

var (
	blockInfoBin = "sols/BlockInfo_sol_BlockInfo.bin"
	blockInfoAbi = "sols/BlockInfo_sol_BlockInfo.abi"
	blockInfoCode []byte
	blockInfoAddress evm.Address
)

func TestBlockInfoSol(t *testing.T) {
	binBytes, err := util.ReadBinFile(blockInfoBin)
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var origin = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	blockInfoCode, blockInfoAddress = deployContract(t, memoryDB, bc, origin, binBytes, "", "nil", 136183)
	callInfo(t, memoryDB, bc, origin, "608060405234801561001057600080fd5b50600436106100935760003560e01c8063ab70fd6911610066578063ab70fd691461013c578063b6baffe31461015a578063d1a82a9d14610178578063df1f29ee146101c2578063f2c9ecd81461020c57610093565b806312065fe014610098578063188ec356146100b657806338cc4831146100d4578063a16963b31461011e575b600080fd5b6100a061022a565b6040518082815260200191505060405180910390f35b6100be610232565b6040518082815260200191505060405180910390f35b6100dc61023a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610126610242565b6040518082815260200191505060405180910390f35b61014461024a565b6040518082815260200191505060405180910390f35b610162610252565b6040518082815260200191505060405180910390f35b61018061025a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6101ca610262565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61021461026a565b6040518082815260200191505060405180910390f35b600047905090565b600042905090565b600030905090565b600045905090565b60003a905090565b600044905090565b600041905090565b600032905090565b60004390509056fea26469706673582212208fb206fed034672fe9f9756ebf8791aafd13cc075ba8e606037265d3a57316c764736f6c63430006000033",183)
}

func callInfo(t *testing.T, db evm.DB, bc evm.Blockchain, caller evm.Address, excepts string, gasCost uint64) {

	var gasQuota uint64 = 10000
	var gas = gasQuota
	output, err := evm.New(bc, db, &evm.Context{
		Input: nil,
		Value: 0,
		Gas: &gas,
	}).Call(caller, blockInfoAddress, blockInfoCode)
	require.NoError(t, err)
	require.Equal(t, excepts, output)
	if gasCost != 0 {
		require.EqualValues(t, gasCost, gasQuota-gas)
	}
}