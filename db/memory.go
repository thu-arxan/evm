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

package db

import (
	"errors"
	"evm"
	"evm/util"
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

// GetAccount is the implementation of interface
func (m *Memory) GetAccount(address evm.Address) evm.Account {
	key := string(address.Bytes())
	if account, ok := m.accounts[key]; ok {
		return account.account
	}
	account := m.accountFunc(address)
	m.accounts[key] = &accountInfo{
		account: account,
	}
	return account
}

// GetStorage is the implementation of interface
func (m *Memory) GetStorage(address evm.Address, key []byte) []byte {
	storageKey := string(util.BytesCombine(address.Bytes(), key))
	if value, ok := m.storages[storageKey]; ok {
		return value
	}
	return nil
}

// NewWriteBatch is the implementation of interface
func (m *Memory) NewWriteBatch() evm.WriteBatch {
	return m
}

// SetStorage is the implementation of interface
func (m *Memory) SetStorage(address evm.Address, key, value []byte) {
	storageKey := string(util.BytesCombine(address.Bytes(), key))
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
