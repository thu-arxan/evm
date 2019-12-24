package evm

import (
	"evm/crypto"
	"evm/rlp"
	"fmt"
)

// This file defines some default funcion if the user do not want to implement it by themself.

// defaultCreateAddress is the default implementation of CreateAddress
func defaultCreateAddress(caller Address, nonce uint64, toAddressFunc func(bytes []byte) Address) Address {
	data, _ := rlp.EncodeToBytes([]interface{}{caller, nonce})
	bytes := crypto.Keccak256(data)[12:]
	return toAddressFunc(bytes)
}

func defaultCreate2Address(caller Address, salt, code []byte, toAddressFunc func(bytes []byte) Address) Address {
	fmt.Printf("%x\n", caller.Bytes())
	bytes := crypto.Keccak256([]byte{0xff}, caller.Bytes(), salt[:], crypto.Keccak256(code))[12:]
	fmt.Printf("%x\n", bytes)
	return toAddressFunc(bytes)
}
