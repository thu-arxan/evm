package evm

// Params defines some params
type Params struct {
	Origin Address
	Caller Address
	Callee Address
	Input  []byte
	Value  uint64
	Gas    *uint64
}
