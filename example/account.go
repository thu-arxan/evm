package example

import "evm"

// Account is account
type Account struct {
	addr    *Address
	code    []byte
	balance uint64
}

// NewAccount is the constructor of Account
func NewAccount(addr *Address) *Account {
	return &Account{
		addr: addr,
	}
}

// SetCode is the implementation of interface
func (a *Account) SetCode(code []byte) {}

// GetAddress is the implementation of interface
func (a *Account) GetAddress() evm.Address {
	return a.addr
}

// GetBalance is the implementation of interface
func (a *Account) GetBalance() uint64 {
	return 100000
}

// GetCode is the implementation of interface
func (a *Account) GetCode() []byte {
	return a.code
}

// GetCodeHash return the hash of account code, please return [32]byte, and return [32]byte{0, ..., 0} if code is empty
func (a *Account) GetCodeHash() []byte {
	var hash = make([]byte, 0)
	return hash
}

// AddBalance is the implementation of interface
func (a *Account) AddBalance(balance uint64) error {
	return nil
}

// SubBalance is the implementation of interface
func (a *Account) SubBalance(balance uint64) error {
	return nil
}
