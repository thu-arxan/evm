package example

import (
	"evm"
	"evm/core"
)

// MemoryDB is a db implementation in memory
type MemoryDB struct {
}

// NewMemoryDB is the constructor of MemoryDB
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{}
}

// Exist is the implementation of interface
func (db *MemoryDB) Exist(address evm.Address) bool {
	return false
}

// GetAccount is the implementation of interface
func (db *MemoryDB) GetAccount(address evm.Address) evm.Account {
	return &Account{}
}

// GetStorage is the implementation of interface
func (db *MemoryDB) GetStorage(address evm.Address, key core.Word256) (value []byte, err error) {
	return nil, nil
}

// SetStorage is the implementation of interface
func (db *MemoryDB) SetStorage(address evm.Address, key core.Word256, value []byte) error {
	return nil
}

// UpdateAccount is the implementation of interface
func (db *MemoryDB) UpdateAccount(account evm.Account) error {
	return nil
}

// RemoveAccount is the implementation of interface
func (db *MemoryDB) RemoveAccount(address evm.Address) error {
	return nil
}

// AddLog is the implementation of interface
func (db *MemoryDB) AddLog(log *evm.Log) {

}
