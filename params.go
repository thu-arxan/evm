package evm

import "evm/core"

// Params defines some params
type Params struct {
	Origin core.Address
	Caller core.Address
	Callee core.Address
	Input  []byte
	Value  uint64
	Gas    *uint64
}
