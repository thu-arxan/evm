package evm

import (
	"evm/core"
	"evm/util"
	"fmt"
)

// Cache cache on a DB.
// It will simulate operate on a db, and sync to db if necessary.
// Note: It's not thread safety because now it will only be used in one thread.
type Cache struct {
	db       DB
	readonly bool
	accounts map[string]*accountInfo
	logs     []*Log
}

type accountInfo struct {
	account Account
	storage map[string][]byte
	removed bool
	updated bool
}

// NewCache is the constructor of Cache
func NewCache(db DB) *Cache {
	return &Cache{
		db:       db,
		accounts: make(map[string]*accountInfo),
	}
}

// Exist return if an account exist
// todo: wrong implementation
func (cache *Cache) Exist(addr Address) bool {
	if util.Contain(cache.accounts, addressToString(addr)) {
		return true
	}
	return cache.db.Exist(addr)
}

// HasSuicide return if an account has suicide
// todo: just return false now
func (cache *Cache) HasSuicide(addr Address) bool {
	return false
}

// GetAccount return the account of address
func (cache *Cache) GetAccount(addr Address) Account {
	return cache.get(addr).account
}

// UpdateAccount set account
func (cache *Cache) UpdateAccount(account Account) error {
	accInfo := cache.get(account.GetAddress())
	if accInfo.removed {
		return fmt.Errorf("UpdateAccount on a removed account: %s", account.GetAddress())
	}
	accInfo.account = account
	accInfo.updated = true
	return nil
}

// Suicide remove an account
func (cache *Cache) Suicide(address Address) error {
	accInfo := cache.get(address)
	if accInfo.removed {
		return fmt.Errorf("RemoveAccount on a removed account: %s", address)
	}
	accInfo.removed = true
	return nil
}

// GetStorage returns the key of an address if exist, else returns an error
func (cache *Cache) GetStorage(address Address, key core.Word256) []byte {
	accInfo := cache.get(address)

	if util.Contain(accInfo.storage, word256ToString(key)) {
		return accInfo.storage[word256ToString(key)]
	}
	value := cache.db.GetStorage(address, key)
	// todo: if we need to deal nil?
	accInfo.storage[word256ToString(key)] = value
	return value
}

// SetStorage set the storage of address
// NOTE: Set value to zero to remove. How should i understand this?
// TODO: Review this
func (cache *Cache) SetStorage(address Address, key core.Word256, value []byte) {
	accInfo := cache.get(address)
	// todo: how to deal account removed
	// if accInfo.removed {
	// 	return fmt.Errorf("SetStorage on a removed account: %s", addressToString(address))
	// }
	accInfo.storage[word256ToString(key)] = value
	accInfo.updated = true
}

// GetNonce return the nonce of account
// todo: implement it in the right way
func (cache *Cache) GetNonce(address Address) uint64 {
	return cache.get(address).account.GetNonce()
}

// AddLog add log
func (cache *Cache) AddLog(log *Log) {
	cache.logs = append(cache.logs, log)
}

// Sync will sync change to db
func (cache *Cache) Sync() {
	wb := cache.db.NewWriteBatch()
	for _, info := range cache.accounts {
		if info.removed {
			wb.Suicide(info.account.GetAddress())
		} else if info.updated {
			wb.UpdateAccount(info.account)
			for key, value := range info.storage {
				wb.SetStorage(info.account.GetAddress(), stringToWord256(key), value)
			}
		}
	}
	for i := range cache.logs {
		wb.AddLog(cache.logs[i])
	}
}

// get the cache accountInfo item creating it if necessary
func (cache *Cache) get(address Address) *accountInfo {
	key := addressToString(address)
	if util.Contain(cache.accounts, key) {
		return cache.accounts[key]
	}
	// Then try to load from db
	// todo: should return error?
	account := cache.db.GetAccount(address)
	// set the account
	cache.accounts[key] = &accountInfo{
		account: account,
		storage: make(map[string][]byte),
		removed: false,
		updated: false,
	}

	return cache.accounts[key]
}

func addressToString(address Address) string {
	return string(address.Bytes())
}

func stringToAddress(s string) Address {
	addr := core.AddressFromBytes([]byte(s))
	return addr
}

func word256ToString(word core.Word256) string {
	return string(word.Bytes())
}

func stringToWord256(s string) core.Word256 {
	return core.BytesToWord256([]byte(s))
}
