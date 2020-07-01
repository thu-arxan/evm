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

package core

import (
	"errors"
	"github.com/thu-arxan/evm/crypto"
	"github.com/thu-arxan/evm/util/math"
)

// Account defines an account structure
type Account struct {
	address  *Address
	code     []byte
	balance  uint64
	nonce    uint64
	suicided bool
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

// Suicide suicide an account
func (a *Account) Suicide() {
	a.suicided = true
}

// HasSuicide return if an account has suicided
func (a *Account) HasSuicide() bool {
	return a.suicided == true
}

// Copy return a copy of an account
// Note: Please reimplement this if account structure changed
func (a *Account) Copy() *Account {
	var account Account
	account.code = CopyBytes(a.code)
	account.address = a.address.Copy()
	account.balance = a.balance
	account.nonce = a.nonce
	account.suicided = a.suicided
	return &account
}

// CopyBytes copy bytes
func CopyBytes(origin []byte) []byte {
	if origin == nil {
		return nil
	}
	var res = make([]byte, len(origin))
	// for i := range origin {
	// 	res[i] = origin[i]
	// }
	copy(res, origin)
	return res
}
