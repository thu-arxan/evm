package evm

import "evm/core"

// Log is the log of evm
type Log struct {
	Address Address
	// list of topics provided by the contract.
	Topics []core.Word256
	// supplied by the contract, usually ABI-encoded
	Data []byte
}
