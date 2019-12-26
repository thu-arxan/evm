package tests

import (
	"errors"
	"evm"
)

// Account is account
type Account struct {
	addr    *Address
	code    []byte
	balance uint64
	nonce   uint64
}

// NewAccount is the constructor of Account
func NewAccount(addr *Address) *Account {
	return &Account{
		addr: addr,
	}
}

// SetCode is the implementation of interface
func (a *Account) SetCode(code []byte) {
	a.code = code
}

// GetAddress is the implementation of interface
func (a *Account) GetAddress() evm.Address {
	return a.addr
}

// GetBalance is the implementation of interface
func (a *Account) GetBalance() uint64 {
	return a.balance
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
// Note: In fact, we should avoid overflow
func (a *Account) AddBalance(balance uint64) error {
	a.balance += balance
	return nil
}

// SubBalance is the implementation of interface
func (a *Account) SubBalance(balance uint64) error {
	if a.balance < balance {
		return errors.New("InsufficientBalance")
	}
	a.balance -= balance
	return nil
}

// GetNonce is the implementation of interface
func (a *Account) GetNonce() uint64 {
	return a.nonce
}

// SetNonce is the implementation of interface
func (a *Account) SetNonce(nonce uint64) {
	a.nonce = nonce
}
