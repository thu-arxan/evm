# Example

This folder provide an example to make use of evm.

To see the documentations of interfaces, refer to README.md.

To get more usage and test about this project, refer to package tests.

## Address

a byte array of length 20 implements Address interface.

## BlockChain

Struct BlockChain implements BlockChain interface.

## Account

Struct Account implements Account interface.

```go
// Account is account
type Account struct {
	addr    *Address
  // code is bytecodes if account is an contract account or nil if not
	code    []byte
  // balance stores the money this account has
	balance uint64
  // nonce is used to privide randomness for this account to create
  // new contract
	nonce   uint64
  // suicide marks if this account is suicided
	suicide bool
}

```

## main.go

An example usage of this project.

```go
// write your solidity files first and use solcjs tool to generate binary and // abi files
// Or use the files we provide
func main() {
  // read binary files and store bytecodes in code
  // these bytecodes are used by evm to get and deploy contract code but
  // not the contract logic itself
  // (a bit confusing but following evm convention)
	code, err := util.ReadBinFile("../sols/output/Balance_sol_Balance.bin")
	if err != nil {
		panic(err)
	}
  // generate dependency
  // including blockchain, gas and database
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var gas uint64
	gas = 10000000
  
  // new a virtual machine to deploy and run the contract
  // code generated above is used as input to evm.create to get contract code
	vm := evm.New(bc, memoryDB, &evm.Context{
		Input: code,
		Value: 0,
		Gas:   &gas,
	})
  
  // pass a random address as caller to create contract
	var caller = example.RandomAddress()
  // now code is actual contract logic and callee is contract address
	code, callee, err := vm.Create(caller)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", code)

	gas = 1000000
  // to call a function inside a contract and pass inputs to it
  // use abi.Pack(abiFiles, function name, inputs...)
	payload, err := abi.Pack("../sols/output/Balance_sol_Balance.abi", "get")
	if err != nil {
		panic(err)
	}
  // use payload generated above as input to a contract
  // and call evm.Call(caller, callee, contract code)
	output, err := evm.New(bc, memoryDB, &evm.Context{
		Input: payload,
		Value: 0,
		Gas:   &gas,
	}).Call(caller, callee, code)
	fmt.Printf("%x\n", output)
}
```

