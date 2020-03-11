//  Copyright 2020 The THU-Arxan Authors
//  This file is part of the evm library.
//
//  The evm library is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Lesser General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  The evm library is distributed in the hope that it will be useful,/
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
//  GNU Lesser General Public License for more details.
//
//  You should have received a copy of the GNU Lesser General Public License
//  along with the evm library. If not, see <http://www.gnu.org/licenses/>.
//

package example

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
	suicide bool
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

// Suicide is the implementation of interface
func (a *Account) Suicide() {
	a.suicide = true
}

// HasSuicide is the implementation of interface
func (a *Account) HasSuicide() bool {
	return a.suicide
}
