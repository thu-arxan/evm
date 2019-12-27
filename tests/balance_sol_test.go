package tests

import (
	"evm"
	"evm/abi"
	"evm/db"
	"evm/example"
	"evm/util"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	balanceAbi     = "sols/Balance_sol_Balance.abi"
	balanceCode    []byte
	balanceAddress evm.Address
)

func TestBalanceSol(t *testing.T) {
	// first create the contract
	binBytes, err := util.ReadBinFile("sols/Balance_sol_Balance.bin")
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var gas uint64 = 1000000
	vm := evm.New(bc, memoryDB, &evm.Context{
		Input: binBytes,
		Value: 0,
		Gas:   &gas,
	})
	var origin = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	code, contractAddress, err := vm.Create(origin)
	require.NoError(t, err)
	var exceptCode = `60806040523480156100115760006000fd5b506004361061005c5760003560e01c80631003e2d21461006257806327ee58a6146100a5578063370158ea146100e857806360fe47b1146101395780636d4ce63c146101805761005c565b60006000fd5b61008f600480360360208110156100795760006000fd5b810190808035906020019092919050505061019e565b6040518082815260200191505060405180910390f35b6100d2600480360360208110156100bc5760006000fd5b81019080803590602001909291905050506101c6565b6040518082815260200191505060405180910390f35b6100f06101ee565b604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b610166600480360360208110156101505760006000fd5b8101908080359060200190929190505050610209565b604051808215151515815260200191505060405180910390f35b610188610225565b6040518082815260200191505060405180910390f35b6000816000600082828250540192505081909090555060006000505490506101c1565b919050565b6000816000600082828250540392505081909090555060006000505490506101e9565b919050565b600060003360006000505481915091509150610205565b9091565b600081600060005081909090555060019050610220565b919050565b60006000600050549050610234565b9056fea26469706673582212206d1f7e72f2d26816fe48ff60de6fa42d7b6fb40fc1603494b8c63221cd3c2c2364736f6c63430006000033`
	require.Equal(t, exceptCode, fmt.Sprintf("%x", code))
	var exceptAddress = `cd234a471b72ba2f1ccf0a70fcaba648a5eecd8d`
	require.Equal(t, exceptAddress, fmt.Sprintf("%x", contractAddress.Bytes()))
	account := memoryDB.GetAccount(contractAddress)
	require.Equal(t, code, account.GetCode())
	require.Equal(t, uint64(1), account.GetNonce())
	balanceAddress = contractAddress
	balanceCode = code
	require.Equal(t, uint64(855600), gas, fmt.Sprintf("except %d while get %d", 855600, gas))
	// then call the contract with get function
	callBalance(t, memoryDB, bc, origin, "get", nil, []string{"10"}, 1096)
	// then set value to 20
	callBalance(t, memoryDB, bc, origin, "set", []string{"20"}, []string{"true"}, 5393)
	// then get
	callBalance(t, memoryDB, bc, origin, "get", nil, []string{"20"}, 1096)
	// then add
	callBalance(t, memoryDB, bc, origin, "add", []string{"10"}, []string{"30"}, 6939)
	// then get
	callBalance(t, memoryDB, bc, origin, "get", nil, []string{"30"}, 1096)
	// info
	callBalance(t, memoryDB, bc, origin, "info", nil, []string{"6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0", "30"}, 1105)
	// define temporary address for testing
	var temporarySender = RandomAddress()
	var temporaryBC = NewBlockchain()
	abi.SetAddressParser(temporarySender.Length(), func(bytes []byte) string {
		return BytesToAddress(bytes).String()
	})
	callBalance(t, memoryDB, temporaryBC, temporarySender, "info", nil, []string{temporarySender.String(), "30"}, 1105)
}

// you can set gasCost to 0 if you do not want to compare gasCost
func callBalance(t *testing.T, db evm.DB, bc evm.Blockchain, caller evm.Address, funcName string, inputs, excepts []string, gasCost uint64) {
	payload, err := abi.GetPayloadBytes(balanceAbi, funcName, inputs)
	require.NoError(t, err)
	var gasQuota uint64 = 10000
	var gas = gasQuota
	output, err := evm.New(bc, db, &evm.Context{
		Input: payload,
		Value: 0,
		Gas:   &gas,
	}).Call(caller, balanceAddress, balanceCode)
	require.NoError(t, err)
	variables, err := abi.Unpacker(balanceAbi, funcName, output)
	require.NoError(t, err)
	require.Len(t, variables, len(excepts))
	for i := range excepts {
		require.Equal(t, excepts[i], variables[i].Value)
	}
	if gasCost != 0 {
		require.EqualValues(t, gasCost, gasQuota-gas)
	}

}
