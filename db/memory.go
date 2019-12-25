package db

import (
	"errors"
	"evm"
	"evm/core"
	"evm/util"
	"fmt"
)

// Memory is a memory db
type Memory struct {
	accounts map[string]evm.Account
	storages map[string][]byte
}

// NewMemory is the constructor of Memory
func NewMemory() *Memory {
	return &Memory{
		accounts: make(map[string]evm.Account),
	}
}

// Exist is the implementation of interface
func (m *Memory) Exist(address evm.Address) bool {
	key := string(address.Bytes())
	return util.Contain(m.accounts, key)
}

// GetAccount is the implementation of interface
func (m *Memory) GetAccount(address evm.Address) evm.Account {
	// return &Account{}
	return nil
}

// GetStorage is the implementation of interface
func (m *Memory) GetStorage(address evm.Address, key core.Word256) (value []byte, err error) {
	storageKey := fmt.Sprintf("%s:%s", address.Bytes(), key.Bytes())
	if util.Contain(m.storages, storageKey) {
		return m.storages[storageKey], nil
	}
	return nil, errors.New("value is not exist")
}

// SetStorage is the implementation of interface
func (m *Memory) SetStorage(address evm.Address, key core.Word256, value []byte) error {
	storageKey := fmt.Sprintf("%s:%s", address.Bytes(), key.Bytes())
	m.storages[storageKey] = value
	return nil
}

// UpdateAccount is the implementation of interface
func (m *Memory) UpdateAccount(account evm.Account) error {
	key := string(account.GetAddress().Bytes())
	m.accounts[key] = account
	return nil
}

// RemoveAccount is the implementation of interface
func (m *Memory) RemoveAccount(address evm.Address) error {
	return nil
}

// GetNonce is the implementation of interface
func (m *Memory) GetNonce(address evm.Address) uint64 {
	return 0
}

// AddLog is the implementation of interface
func (m *Memory) AddLog(log *evm.Log) {

}
