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

package evm

// This file defines some kinds of interfaces

// Account describe what function that account should provide
type Account interface {
	GetAddress() Address
	GetBalance() uint64
	AddBalance(balance uint64) error
	SubBalance(balance uint64) error
	GetCode() []byte
	SetCode(code []byte)
	// GetCodeHash return the hash of account code, please return [32]byte,
	// and return [32]byte{0, ..., 0} if code is empty
	GetCodeHash() []byte
	GetNonce() uint64
	SetNonce(nonce uint64)
	// Suicide will suicide an account
	Suicide()
	HasSuicide() bool
}

// Address describe what functions that an Address implementation should provide
type Address interface {
	// It would be better if length = 32
	// 1. Add zero in left if length < 32
	// 2. Remove left byte if length > 32(however, this may be harm)
	Bytes() []byte
}

// DB describe what function that db should provide to support the evm
type DB interface {
	// Exist return if the account exist
	// Note: if account is suicided, return true
	Exist(address Address) bool
	// GetStorage return a default account if unexist
	GetAccount(address Address) Account
	// Note: GetStorage return nil if key is not exist
	GetStorage(address Address, key []byte) (value []byte)
	NewWriteBatch() WriteBatch
}

// WriteBatch define a batch which support some write operations
type WriteBatch interface {
	SetStorage(address Address, key []byte, value []byte)
	// Note: db should delete all storages if an account suicide
	UpdateAccount(account Account) error
	AddLog(log *Log)
}

// Blockchain describe what function that blockchain system shoudld provide to support the evm
type Blockchain interface {
	// GetBlockHash return ZeroWord256 if num > 256 or num > max block height
	GetBlockHash(num uint64) []byte
	// CreateAddress will be called by CREATE Opcode
	CreateAddress(caller Address, nonce uint64) Address
	// Create2Address will be called by CREATE2 Opcode
	Create2Address(caller Address, salt, code []byte) Address
	// Note: NewAccount will create a default account in Blockchain service,
	// but please do not append the account into db here
	NewAccount(address Address) Account
	// BytesToAddress provide a way convert bytes(normally [32]byte) to Address
	BytesToAddress(bytes []byte) Address
}
