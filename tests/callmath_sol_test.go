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
	callMathAbi     = "sols/CallMath_sol_Call.abi"
	callMathBin     = "sols/CallMath_sol_Call.bin"
	callMathCode    []byte
	callMathAddress evm.Address
)

func TestCallMath(t *testing.T) {
	var err error
	binBytes, err := util.ReadBinFile(mathBin)
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var origin = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	var exceptCode = `60806040523480156100115760006000fd5b50600436106100305760003560e01c8063bdf118791461003657610030565b60006000fd5b61003e610054565b6040518082815260200191505060405180910390f35b600060011515600260009054906101000a900460ff161515141561008b576001600050546001600050540260016000508190909055505b60007f61000000000000000000000000000000000000000000000000000000000000009050600060046000505460036020811015156100c657fe5b1a60f81b9050807effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191610600260006101000a81548160ff0219169083151502179055506002600060005054026000600050819090905550600660006000505481151561015157fe5b05600060005081909090555060046000600050540a6000600050819090905550600460036000505401600060005054600082121561018b57fe5b901b600060005081909090555060016003600050540160006000505460008212156101b257fe5b901d60006000508190909055506003600050546000600050548115156101d457fe5b0760006000508190909055506003600050548015156101ef57fe5b60016000505460016000505408600060005081909090555060036000505480151561021657fe5b600160005054600160005054096000600050819090905550600060005054600160005081909090555061012c600160005054111561026c57600360016000505481151561025f57fe5b0660016000508190909055505b6000600050546000600050541315156102f357600060005054600160005054166001600050819090905550600160005054196001600050819090905550600060005054600160005054176001600050819090905550600360005054600160008282825054019250508190909055506003600050546001600050541860016000508190909055505b600160005054925050506103045650505b9056fea264697066735822122035b640f0acd6978a23cfef08594e98f37eefc72762fb9624986d5cdbe7e0d7bb64736f6c63430006000033`
	mathCode, mathAddress = deployContract(t, memoryDB, bc, origin, binBytes, "cd234a471b72ba2f1ccf0a70fcaba648a5eecd8d", exceptCode, 246938)
	// then we deploy call math contract
	binBytes, err = util.ReadBinFile(callMathBin)
	exceptCode = `608060405234801560105760006000fd5b5060043610602c5760003560e01c8063d9c6717414603257602c565b60006000fd5b6038604e565b6040518082815260200191505060405180910390f35b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663bdf118796040518163ffffffff1660e01b8152600401602060405180830381600087803b15801560ba5760006000fd5b505af115801560ce573d600060003e3d6000fd5b505050506040513d602081101560e45760006000fd5b8101908080519060200190929190505050905060fb565b9056fea26469706673582212202422e8bee4ee2a1ab327cff2975e43e3eaa268c19bf66911de89729f7dca557064736f6c63430006000033`
	callMathCode, callMathAddress = deployContract(t, memoryDB, bc, origin, binBytes, "343c43a37d37dff08ae8c4a11544c718abb4fcf8", exceptCode, 104293)
	// then we call
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	call(t, memoryDB, bc, origin, callMathAddress, callMathAbi, "callMath", nil, []string{"1"}, 56605, 30000)
}

func call(t *testing.T, db evm.DB, bc evm.Blockchain, caller, contract evm.Address, abiFile, funcName string, inputs, excepts []string, gasCost, refund uint64) {
	payload, err := abi.GetPayloadBytes(abiFile, funcName, inputs)
	require.NoError(t, err)
	var gasQuota uint64 = 100000
	var gas = gasQuota
	fmt.Printf("Payload is %x\n", payload)
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
