package core

import (
	"errors"
	"evm/crypto"
	"evm/util/math"
)

// Account defines an account structure
type Account struct {
	address *Address
	code    []byte
	balance uint64
	nonce   uint64
}

// NewAccount is the constructor of Account
func NewAccount(address *Address) *Account {
	return &Account{
		address: address,
	}
}

// GetAddress return the address
func (a *Account) GetAddress() *Address {
	return a.address
}

// GetBalance return the balance
func (a *Account) GetBalance() uint64 {
	return a.balance
}

// AddBalance add balance
func (a *Account) AddBalance(balance uint64) error {
	sum, overflow := math.SafeAdd(a.balance, balance)
	if overflow {
		return errors.New("overflow")
	}
	a.balance = sum
	return nil
}

// SubBalance sub balance
func (a *Account) SubBalance(balance uint64) error {
	sub, overflow := math.SafeSub(a.balance, balance)
	if overflow {
		return errors.New("insufficient balance")
	}
	a.balance = sub
	return nil
}

// GetCode return the code
func (a *Account) GetCode() []byte {
	return a.code
}

// SetCode set code
func (a *Account) SetCode(code []byte) {
	a.code = code
}

// GetCodeHash return the hash of code
func (a *Account) GetCodeHash() []byte {
	return crypto.Keccak256(a.code)
}

// GetNonce return the nonce of account
func (a *Account) GetNonce() uint64 {
	return a.nonce
}

// SetNonce set the nonce of account
func (a *Account) SetNonce(nonce uint64) {
	a.nonce = nonce
}
