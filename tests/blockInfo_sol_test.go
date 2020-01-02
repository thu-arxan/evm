package tests

import (
	"evm"
	"evm/abi"
	"evm/db"
	"evm/example"
	"evm/util"

	"testing"
	"github.com/stretchr/testify/require"
)

var (
	blockInfoBin = "sols/BlockInfo_sol_SimpleBlock.bin"
	blockInfoAbi = "sols/BlockInfo_sol_SimpleBlock.abi"
	blockInfoCode []byte
	blockInfoAddress evm.Address
)

func TestBlockInfoSol(t *testing.T) {
	binBytes, err := util.ReadBinFile(blockInfoBin)
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var origin = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
//	var exceptCode = `60806040523480156100115760006000fd5b50610017565b6102db806100266000396000f3fe60806040523480156100115760006000fd5b50600436106100985760003560e01c8063ab70fd6911610067578063ab70fd6914610142578063b6baffe314610160578063d1a82a9d1461017e578063df1f29ee146101c8578063f2c9ecd81461021257610098565b806312065fe01461009e578063188ec356146100bc57806338cc4831146100da578063a16963b31461012457610098565b60006000fd5b6100a6610230565b6040518082815260200191505060405180910390f35b6100c461023d565b6040518082815260200191505060405180910390f35b6100e261024a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61012c610257565b6040518082815260200191505060405180910390f35b61014a610264565b6040518082815260200191505060405180910390f35b610168610271565b6040518082815260200191505060405180910390f35b61018661027e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6101d061028b565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61021a610298565b6040518082815260200191505060405180910390f35b600047905061023a565b90565b6000429050610247565b90565b6000309050610254565b90565b6000459050610261565b90565b60003a905061026e565b90565b600044905061027b565b90565b6000419050610288565b90565b6000329050610295565b90565b60004390506102a2565b9056fea264697066735822122015c6e882c3fdc1443e57fbd751c159a2579310479893e31b21797f9e2579ce4b64736f6c63430006000033`
	var exceptAddress = `cd234a471b72ba2f1ccf0a70fcaba648a5eecd8d`
	blockInfoCode, blockInfoAddress = deployContract(t, memoryDB, bc, origin, binBytes, exceptAddress, "", 0)
	callInfo(t, memoryDB, bc, origin, "getAddress", nil, nil, 281)
	callInfo(t, memoryDB, bc, origin, "getBalance", nil, nil, 228)
	callInfo(t, memoryDB, bc, origin, "getOrigin", nil, nil, 302)
	callInfo(t, memoryDB, bc, origin, "getGasprice", nil, nil, 224)
	callInfo(t, memoryDB, bc, origin, "getCoinbase", nil, nil, 280)
	callInfo(t, memoryDB, bc, origin, "getTimestamp", nil, nil, 247)
	callInfo(t, memoryDB, bc, origin, "getNumber", nil, nil, 312)
	callInfo(t, memoryDB, bc, origin, "getDifficulty", nil, nil, 246)
	callInfo(t, memoryDB, bc, origin, "getGaslimit", nil, nil, 291)
	callInfo(t, memoryDB, bc, origin, "getChainID", nil, nil, 304)

}

func callInfo(t *testing.T, db evm.DB, bc evm.Blockchain, caller evm.Address,  funcName string, inputs, excepts []string, gasCost uint64) {
	payload, err := abi.GetPayloadBytes(blockInfoAbi, funcName, inputs)
	require.NoError(t, err)
	var gasQuota uint64 = 10000
	var gas = gasQuota
	output, err := evm.New(bc, db, &evm.Context{
		Input: payload,
		Value: 0,
		Gas: &gas,
	}).Call(caller, blockInfoAddress, blockInfoCode)
	require.NoError(t, err)
	if excepts != nil {
		variables, err := abi.Unpacker(balanceAbi, funcName, output)
		require.NoError(t, err)
		require.Len(t, variables, len(excepts))
		for i := range excepts {
			require.Equal(t, excepts[i], variables[i].Value)
		}
	}
	if gasCost != 0 {
		require.EqualValues(t, gasCost, gasQuota-gas)
	}
	t.Log(output)
	t.Log(gasQuota - gas)
}
