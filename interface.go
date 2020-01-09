package evm

import "evm/core"

// This file defines some kinds of interfaces

// Account describe what function that account should provide
type Account interface {
	GetAddress() Address
	GetBalance() uint64
	AddBalance(balance uint64) error
	SubBalance(balance uint64) error
	GetCode() []byte
	SetCode(code []byte)
	// GetCodeHash return the hash of account code, please return [32]byte, and return [32]byte{0, ..., 0} if code is empty
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
	GetStorage(address Address, key core.Word256) (value []byte)
	NewWriteBatch() WriteBatch
}

// WriteBatch define a batch which support some write operations
type WriteBatch interface {
	SetStorage(address Address, key core.Word256, value []byte)
	// Note: db should delete all storages if an account suicide
	UpdateAccount(account Account) error
	AddLog(log *Log)
}

// Blockchain describe what function that blockchain system shoudld provide to support the evm
type Blockchain interface {
	GetBlockHash(num uint64) ([]byte, error)
	// CreateAddress will be called by CREATE Opcode
	CreateAddress(caller Address, nonce uint64) Address
	// Create2Address will be called by CREATE2 Opcode
	Create2Address(caller Address, salt, code []byte) Address
	// Note: NewAccount will create a default account in Blockchain service, but please do not append the account into db here
	NewAccount(address Address) Account
	// BytesToAddress provide a way convert bytes(normally [32]byte) to Address
	BytesToAddress(bytes []byte) Address
}

// emptyAccount contain nothing
type emptyAccount struct{}

func (account emptyAccount) SetCode(code []byte)             {}
func (account emptyAccount) GetAddress() Address             { return nil }
func (account emptyAccount) GetBalance() uint64              { return 0 }
func (account emptyAccount) GetEVMCode() []byte              { return nil }
func (account emptyAccount) GetCode() []byte                 { return nil }
func (account emptyAccount) GetCodeHash() []byte             { return nil }
func (account emptyAccount) SubBalance(balance uint64) error { return nil }
func (account emptyAccount) AddBalance(balance uint64) error { return nil }
