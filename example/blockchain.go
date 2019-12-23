package example

import "evm"

// Blockchain is the implementation of blockchain
type Blockchain struct {
}

// NewBlockchain is the constructor of Blockchain
func NewBlockchain() *Blockchain {
	return &Blockchain{}
}

// GetBlockHash is the implementation of interface
func (bc *Blockchain) GetBlockHash(num uint64) ([]byte, error) {
	return nil, nil
}

// CreateAddress is the implementation of interface
func (bc *Blockchain) CreateAddress(caller evm.Address, nonce uint64) evm.Address {
	return nil
}

// Create2Address is the implementation of interface
func (bc *Blockchain) Create2Address(caller evm.Address, salt, code []byte) evm.Address {
	return nil
}

// NewAccount is the implementation of interface
func (bc *Blockchain) NewAccount(address evm.Address) evm.Account {
	return nil
}

// BytesToAddress is the implementation of interface
func (bc *Blockchain) BytesToAddress(bytes []byte) evm.Address {
	return nil
}
