package evm

// Context defines some context
type Context struct {
	Input []byte
	Value uint64
	Gas   *uint64

	BlockHeight uint64
	BlockTime   int64
	Difficulty  uint64
	GasLimit    uint64
	GasPrice    uint64
	CoinBase    []byte
}
