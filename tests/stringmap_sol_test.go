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
	stringmapAbi     = "sols/StringMap_sol_StringMap.abi"
	stringmapBin     = "sols/StringMap_sol_StringMap.bin"
	stringmapCode    []byte
	stringmapAddress evm.Address
)

func TestStringMapSol(t *testing.T) {
	var err error
	binBytes, err := util.ReadBinFile(stringmapBin)
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var origin = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	var exceptCode = `60806040523480156100115760006000fd5b506004361061005c5760003560e01c806323814fc5146100625780632aac40391461019f5780636833d54f146102dc57806380599e4b146103b7578063ebdf86ca1461047a5761005c565b60006000fd5b610123600480360360208110156100795760006000fd5b81019080803590602001906401000000008111156100975760006000fd5b8201836020820111156100aa5760006000fd5b803590602001918460018302840111640100000000831117156100cd5760006000fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509090919290909192905050506105db565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101645780820151818401525b602081019050610148565b50505050905090810190601f1680156101915780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610260600480360360208110156101b65760006000fd5b81019080803590602001906401000000008111156101d45760006000fd5b8201836020820111156101e75760006000fd5b8035906020019184600183028401116401000000008311171561020a5760006000fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509090919290909192905050506106a4565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156102a15780820151818401525b602081019050610285565b50505050905090810190601f1680156102ce5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61039d600480360360208110156102f35760006000fd5b81019080803590602001906401000000008111156103115760006000fd5b8201836020820111156103245760006000fd5b803590602001918460018302840111640100000000831117156103475760006000fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509090919290909192905050506107bd565b604051808215151515815260200191505060405180910390f35b610478600480360360208110156103ce5760006000fd5b81019080803590602001906401000000008111156103ec5760006000fd5b8201836020820111156103ff5760006000fd5b803590602001918460018302840111640100000000831117156104225760006000fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509090919290909192905050506108e4565b005b6105d9600480360360408110156104915760006000fd5b81019080803590602001906401000000008111156104af5760006000fd5b8201836020820111156104c25760006000fd5b803590602001918460018302840111640100000000831117156104e55760006000fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509090919290909192908035906020019064010000000081111561054d5760006000fd5b8201836020820111156105605760006000fd5b803590602001918460018302840111640100000000831117156105835760006000fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050909091929090919290505050610963565b005b6000600050818051602081018201805184825260208301602085012081835280955050505050506000915090508054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561069c5780601f106106715761010080835404028352916020019161069c565b820191906000526020600020905b81548152906001019060200180831161067f57829003601f168201915b505050505081565b60606000600050826040518082805190602001908083835b6020831015156106e257805182525b6020820191506020810190506020830392506106bc565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390206000508054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107ac5780601f10610781576101008083540402835291602001916107ac565b820191906000526020600020905b81548152906001019060200180831161078f57829003601f168201915b505050505090506107b8565b919050565b6000604051602001808050600001905060405160208183030381529060405280519060200120600019166000600050836040518082805190602001908083835b60208310151561082357805182525b6020820191506020810190506020830392506107fd565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902060005060405160200180828054600181600116156101000203166002900480156108b75780601f106108955761010080835404028352918201916108b7565b820191906000526020600020905b8154815290600101906020018083116108a3575b50509150506040516020818303038152906040528051906020012060001916141590506108df565b919050565b6000600050816040518082805190602001908083835b60208310151561092057805182525b6020820191506020810190506020830392506108fa565b6001836020036101000a0380198251168184511680821785525050505050509050019150509081526020016040518091039020600061095f91906109ef565b5b50565b806000600050836040518082805190602001908083835b6020831015156109a057805182525b60208201915060208101905060208303925061097a565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902060005090805190602001906109e9929190610a37565b505b5050565b50805460018160011615610100020316600290046000825580601f10610a155750610a34565b601f016020900490600052602060002090810190610a339190610abc565b5b50565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610a7857805160ff1916838001178555610aab565b82800160010185558215610aab579182015b82811115610aaa5782518260005090905591602001919060010190610a8a565b5b509050610ab89190610abc565b5090565b610ae49190610ac6565b80821115610ae05760008181506000905550600101610ac6565b5090565b9056fea2646970667358221220a66dfe680ee0de40098acb9b6c230c10b5ef6fb39a6425ab6824f0cc7199036964736f6c63430006000033`
	var exceptAddress = `cd234a471b72ba2f1ccf0a70fcaba648a5eecd8d`
	stringmapCode, stringmapAddress = deployContract(t, memoryDB, bc, origin, binBytes, exceptAddress, exceptCode, 569612)
	callWithPayload(t, memoryDB, bc, origin, stringmapAddress, mustPack(stringmapAbi, "add", "aaa", "bbb"), 0, 0)
	result := callWithPayload(t, memoryDB, bc, origin, stringmapAddress, mustPack(stringmapAbi, "getByKey", "aaa"), 0, 0)
	require.EqualValues(t, []string{"bbb"}, mustUnpack(stringmapAbi, "getByKey", result))
}
