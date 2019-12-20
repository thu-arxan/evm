package evm

// Context defines some context
type Context struct {
	Origin Address
	Caller Address
	Callee Address
	Input  []byte
	Value  uint64
	Gas    *uint64
}
