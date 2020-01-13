package tests

import (
	"evm"
	"evm/db"
	"evm/example"
	"evm/util"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	moneyAbi     = "sols/Money_sol_Money.abi"
	moneyBin     = "sols/Money_sol_Money.bin"
	moneyCode    []byte
	moneyAddress evm.Address
)

func TestMoneySol(t *testing.T) {
	var err error
	binBytes, err := util.ReadBinFile(moneyBin)
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var user = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	require.NoError(t, memoryDB.InitBalance(user, 100))
	var exceptCode = `6080604052600436106100435760003560e01c80634f2be91f146100525780636bdebcc91461005c5780636d4ce63c14610074578063a9059cbb146100a05761004c565b3661004c575b5b005b60006000fd5b61005a6100fd565b005b3480156100695760006000fd5b50610072610100565b005b3480156100815760006000fd5b5061008a61011b565b6040518082815260200191505060405180910390f35b3480156100ad5760006000fd5b506100fb600480360360408110156100c55760006000fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610128565b005b5b565b3373ffffffffffffffffffffffffffffffffffffffff16ff5b565b6000479050610125565b90565b8173ffffffffffffffffffffffffffffffffffffffff166108fc60024781151561014e57fe5b0483019081150290604051600060405180830381858888f1935050505015801561017d573d600060003e3d6000fd5b505b505056fea26469706673582212205e3b6c0a09f50cc0eeefe5374fa3f552169d2f9a99e3a525d262645bdd4ef6d264736f6c63430006000033`
	moneyCode, moneyAddress = deployContractWithValue(t, memoryDB, bc, user, binBytes, 10, 88325)
	require.Equal(t, exceptCode, fmt.Sprintf("%x", moneyCode))
	// then call get
	result := callWithPayload(t, memoryDB, bc, user, moneyAddress, mustPack(moneyAbi, "get"), 249, 0)
	require.EqualValues(t, []string{"10"}, mustUnpack(moneyAbi, "get", result))
	// add value 10 and get will return 20
	callWithPayloadAndValue(t, memoryDB, bc, user, moneyAddress, mustPack(moneyAbi, "add"), 10, 99, 0)
	result = callWithPayload(t, memoryDB, bc, user, moneyAddress, mustPack(moneyAbi, "get"), 249, 0)
	require.EqualValues(t, []string{"20"}, mustUnpack(moneyAbi, "get", result))
	// then we will call withpayload empty and contract money is 30, while user is 70
	callWithPayloadAndValue(t, memoryDB, bc, user, moneyAddress, nil, 10, 57, 0)
	result = callWithPayload(t, memoryDB, bc, user, moneyAddress, mustPack(moneyAbi, "get"), 249, 0)
	require.EqualValues(t, []string{"30"}, mustUnpack(moneyAbi, "get", result))
	require.EqualValues(t, 70, memoryDB.GetAccount(user).GetBalance())
	var someone = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf1")
	// contract has 30, user has 70
	// then someone will has 5 + 15 = 20, contract has 10
	result = callWithPayload(t, memoryDB, bc, user, moneyAddress, mustPack(moneyAbi, "transfer", "6ac7ea33f8831ea9dcc53393aaa88b25a785dbf1", "5"), 32851, 0)
	require.EqualValues(t, 20, memoryDB.GetAccount(someone).GetBalance())
	// then we call transfer to
	// destory
	result = callWithPayload(t, memoryDB, bc, user, moneyAddress, mustPack(moneyAbi, "destory"), 5143, 24000)
	require.Len(t, result, 0)
	require.EqualValues(t, 80, memoryDB.GetAccount(user).GetBalance(), fmt.Sprintf("except 80 other than %d", memoryDB.GetAccount(user).GetBalance()))
}
