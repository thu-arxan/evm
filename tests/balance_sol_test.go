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

func TestBalanceSol(t *testing.T) {
	// first create the contract
	binBytes, err := util.ReadBinFile("sols/Balance_sol_Balance.bin")
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var gas uint64 = 10000
	vm := evm.New(bc, memoryDB, &evm.Context{
		Input: binBytes,
		Value: 0,
		Gas:   &gas,
	})
	var origin = example.HexToAddress("0x6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	code, contractAddress, err := vm.Create(origin)
	require.NoError(t, err)
	var exceptCode = `60806040523480156100115760006000fd5b506004361061005c5760003560e01c80631003e2d21461006257806327ee58a6146100a5578063370158ea146100e857806360fe47b1146101395780636d4ce63c146101805761005c565b60006000fd5b61008f600480360360208110156100795760006000fd5b810190808035906020019092919050505061019e565b6040518082815260200191505060405180910390f35b6100d2600480360360208110156100bc5760006000fd5b81019080803590602001909291905050506101c6565b6040518082815260200191505060405180910390f35b6100f06101ee565b604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b610166600480360360208110156101505760006000fd5b8101908080359060200190929190505050610209565b604051808215151515815260200191505060405180910390f35b610188610225565b6040518082815260200191505060405180910390f35b6000816000600082828250540192505081909090555060006000505490506101c1565b919050565b6000816000600082828250540392505081909090555060006000505490506101e9565b919050565b600060003360006000505481915091509150610205565b9091565b600081600060005081909090555060019050610220565b919050565b60006000600050549050610234565b9056fea26469706673582212206d1f7e72f2d26816fe48ff60de6fa42d7b6fb40fc1603494b8c63221cd3c2c2364736f6c63430006000033`
	require.Equal(t, exceptCode, fmt.Sprintf("%x", code))
	var exceptAddress = `cd234a471b72ba2f1ccf0a70fcaba648a5eecd8d`
	require.Equal(t, exceptAddress, fmt.Sprintf("%x", contractAddress.Bytes()))
	account := memoryDB.GetAccount(contractAddress)
	require.Equal(t, code, account.GetCode())
	require.Equal(t, uint64(1), account.GetNonce())
	// then call the contract
	payload, err := abi.GetPayloadBytes("sols/Balance_sol_Balance.abi", "get", nil)
	require.NoError(t, err)
	fmt.Printf("%x\n", payload)
	vm = evm.New(bc, memoryDB, &evm.Context{
		Input: payload,
		Value: 0,
		Gas:   &gas,
	})
	output, err := vm.Call(origin, contractAddress, code)
	require.NoError(t, err)
	fmt.Printf("%x\n", output)
}
