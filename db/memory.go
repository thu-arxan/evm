package db

import (
	"errors"
	"evm"
	"evm/core"
	"evm/util"
	"fmt"
	"strings"
)

// Memory is a memory db
type Memory struct {
	accounts map[string]*accountInfo
	storages map[string][]byte
	logs     []*evm.Log

	accountFunc func(address evm.Address) evm.Account
}

type accountInfo struct {
	account evm.Account
	removed bool
}

// NewMemory is the constructor of Memory
func NewMemory(accountFunc func(address evm.Address) evm.Account) *Memory {
	return &Memory{
		accounts:    make(map[string]*accountInfo),
		storages:    make(map[string][]byte),
		logs:        make([]*evm.Log, 0),
		accountFunc: accountFunc,
	}
}

// InitBalance init balance for testting, and return error if the account exsit aleardy
func (m *Memory) InitBalance(address evm.Address, balance uint64) error {
	key := string(address.Bytes())
	if util.Contain(m.accounts, key) {
		return errors.New("initial balance on an exist account")
	}
	account := m.accountFunc(address)
	account.AddBalance(balance)
	m.accounts[key] = &accountInfo{
		account: account,
	}
	return nil
}

// Exist is the implementation of interface
func (m *Memory) Exist(address evm.Address) bool {
	key := string(address.Bytes())
	return util.Contain(m.accounts, key)
}

// HasSuicide is the implementation of interface
// todo: it just return false now.
func (m *Memory) HasSuicide(address evm.Address) bool {
	return false
}

// GetAccount is the implementation of interface
func (m *Memory) GetAccount(address evm.Address) evm.Account {
	key := string(address.Bytes())
	if util.Contain(m.accounts, key) {
		return m.accounts[key].account
	}
	account := m.accountFunc(address)
	m.accounts[key] = &accountInfo{
		account: account,
	}
	return account
}

// GetStorage is the implementation of interface
func (m *Memory) GetStorage(address evm.Address, key core.Word256) []byte {
	storageKey := fmt.Sprintf("%s:%s", address.Bytes(), key.Bytes())
	if util.Contain(m.storages, storageKey) {
		return m.storages[storageKey]
	}
	return nil
}

// NewWriteBatch is the implementation of interface
func (m *Memory) NewWriteBatch() evm.WriteBatch {
	return m
}

// SetStorage is the implementation of interface
func (m *Memory) SetStorage(address evm.Address, key core.Word256, value []byte) {
	storageKey := fmt.Sprintf("%s:%s", address.Bytes(), key.Bytes())
	m.storages[storageKey] = value
}

// UpdateAccount is the implementation of interface
func (m *Memory) UpdateAccount(account evm.Account) error {
	key := string(account.GetAddress().Bytes())
	if util.Contain(m.accounts, key) {
		if m.accounts[key].removed {
			return errors.New("can not update on removed account")
		}
	}
	m.accounts[key] = &accountInfo{
		account: account,
	}
	return nil
}

// Suicide is the implementation of interface
// TODO: What if an acount is not exist?
func (m *Memory) Suicide(address evm.Address) error {
	key := string(address.Bytes())
	if util.Contain(m.accounts, key) {
		m.accounts[key].removed = true
	}
	// remove all storages
	for storageKey := range m.storages {
		if strings.HasPrefix(storageKey, key) {
			delete(m.storages, storageKey)
		}
	}
	return nil
}

// AddLog is the implementation of interface
func (m *Memory) AddLog(log *evm.Log) {
	// Note: We should set some infos like txIndex, blockHash and etc.
	// We just set index as example.
	log.Index = uint(len(m.logs))
	m.logs = append(m.logs, log)
}

// GetLog return logs
func (m *Memory) GetLog() []*evm.Log {
	return m.logs
}
