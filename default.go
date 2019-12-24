package evm

import (
	"evm/crypto"
	"evm/rlp"
)

// This file defines some default funcion if the user do not want to implement it by themself.

// defaultCreateAddress is the default implementation of CreateAddress
func defaultCreateAddress(caller Address, nonce uint64, toAddressFunc func(bytes []byte) Address) Address {
	data, _ := rlp.EncodeToBytes([]interface{}{caller, nonce})
	bytes := crypto.Keccak256(data)[12:]
	return toAddressFunc(bytes)
}
