package evm

import "evm/core"

// This file defines some kinds of interfaces

// Account describe what function that account should provide
type Account interface {
	SetCode(code []byte)
	GetAddress() Address
	GetBalance() uint64
	GetEVMCode() []byte
	// todo: what is the difference between GetEVMCode ans GetCode
	GetCode() []byte
	GetCodeHash() []byte
	AddBalance(balance uint64) error
}

// Address describe what functions that an Address implementation should provide
type Address interface {
	Word256() core.Word256
}

// DB describe what function that db should provide to support an evm
type DB interface {
	GetAccount(address Address) (Account, error)
	GetStorage(address Address, key core.Word256) (value []byte, err error)
	SetStorage(address Address, key core.Word256, value []byte) error
	UpdateAccount(account Account) error
	// Remove the account at address
	RemoveAccount(address Address) error
}
