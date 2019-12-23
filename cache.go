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
	accounts map[string]*accountInfo
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
func (cache *Cache) Exist(addr Address) bool {
	if util.Contain(cache.accounts, string(addr.Bytes())) {
		return true
	}
	return cache.db.Exist(addr)
}

// GetAccount return the account of address
func (cache *Cache) GetAccount(addr Address) Account {
	return cache.get(addr).account
}

// SetAccount set account
func (cache *Cache) SetAccount(account Account) error {
	accInfo := cache.get(account.GetAddress())
	if accInfo.removed {
		return fmt.Errorf("UpdateAccount on a removed account: %s", account.GetAddress())
	}
	accInfo.account = account
	accInfo.updated = true
	return nil
}

// RemoveAccount remove an account
func (cache *Cache) RemoveAccount(address Address) error {
	accInfo := cache.get(address)
	if accInfo.removed {
		return fmt.Errorf("RemoveAccount on a removed account: %s", address)
	}
	accInfo.removed = true
	return nil
}

// GetStorage returns the key of an address if exist, else returns an error
func (cache *Cache) GetStorage(address Address, key core.Word256) ([]byte, error) {
	// fmt.Printf("GetStorage of address %s and key %b\n", address.String(), key)
	accInfo := cache.get(address)

	if util.Contain(accInfo.storage, word256ToString(key)) {
		return accInfo.storage[word256ToString(key)], nil
	}
	value, err := cache.db.GetStorage(address, key)
	if err != nil {
		return core.Zero256.Bytes(), err
	}
	accInfo.storage[word256ToString(key)] = value
	return value, nil
}

// SetStorage set the storage of address
// NOTE: Set value to zero to remove. How should i understand this?
func (cache *Cache) SetStorage(address Address, key core.Word256, value []byte) error {
	// fmt.Printf("!!!Set storage %s at key %b and value is %b\n", address.String(), key, value)
	accInfo := cache.get(address)
	if accInfo.removed {
		return fmt.Errorf("SetStorage on a removed account: %s", string(address.Bytes()))
	}
	accInfo.storage[word256ToString(key)] = value
	accInfo.updated = true
	return nil
}

// Sync sync changes to db
// If the sync return an error, it may cause something wrong, so it should be
// deal with by the developer.
// Also, this function may deal with the address and key in an order, so this
// function should be rethink if necessary.
// TODO: Sync should panic rather than return an error
// func (cache *Cache) Sync(wb db.WriteBatch) error {
// 	var err error
// 	for address, account := range cache.accounts {
// 		if account.removed {
// 			if err = wb.RemoveAccount(stringToAddress(address)); err != nil {
// 				return err
// 			}
// 		} else if account.updated {
// 			// err = wb.SetAccount(account.account)
// 			// if err != nil {
// 			// 	return err
// 			// }
// 			// for key, value := range account.storage {
// 			// 	if err = wb.SetStorage(stringToAddress(address), stringToWord256(key), value); err != nil {
// 			// 		return err
// 			// 	}
// 			// }
// 		}
// 	}
// 	return nil
// }

// get the cache accountInfo item creating it if necessary
func (cache *Cache) get(address Address) *accountInfo {
	if util.Contain(cache.accounts, string(address.Bytes())) {
		return cache.accounts[string(address.Bytes())]
	}
	// Then try to load from db
	// todo: should return error?
	account := cache.db.GetAccount(address)
	// set the account
	cache.accounts[string(address.Bytes())] = &accountInfo{
		account: account,
		storage: make(map[string][]byte),
		removed: false,
		updated: false,
	}

	return cache.accounts[string(address.Bytes())]
}

func addressToString(address Address) string {
	return string(address.Bytes())
}

func stringToAddress(s string) Address {
	addr, _ := core.AddressFromBytes([]byte(s))
	return addr
}

func word256ToString(word core.Word256) string {
	return string(word.Bytes())
}

func stringToWord256(s string) core.Word256 {
	word, _ := core.BytesToWord256([]byte(s))
	return word
}
