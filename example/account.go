package example

import "evm"

// Account is account
type Account struct {
}

// SetCode is the implementation of interface
func (a *Account) SetCode(code []byte) {}

// GetAddress is the implementation of interface
func (a *Account) GetAddress() evm.Address {
	return &Address{}
}

// GetBalance is the implementation of interface
func (a *Account) GetBalance() uint64 {
	return 0
}

// GetCode is the implementation of interface
func (a *Account) GetCode() []byte {
	return nil
}

// GetCodeHash return the hash of account code, please return [32]byte, and return [32]byte{0, ..., 0} if code is empty
func (a *Account) GetCodeHash() []byte {
	return nil
}

// AddBalance is the implementation of interface
func (a *Account) AddBalance(balance uint64) error {
	return nil
}

// SubBalance is the implementation of interface
func (a *Account) SubBalance(balance uint64) error {
	return nil
}
