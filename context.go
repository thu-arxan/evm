package evm

import (
	"evm/core"
)

// Context provide a context to run a contract on the evm
type Context struct {
	ChannelID string
	// Number is the number of the block
	Number uint64
	// BlockHash is the hash of the block
	BlockHash core.Word256
	// BlockTime is the time of the block
	BlockTime int64
	// GasLimit limit the use of gas, now is useless
	GasLimit uint64
	// CoinBase, set it to zero
	CoinBase core.Word256
	// diffculty is zero
	Diffculty uint64
}

// GetBlockHash provide a way to get block hash
// todo: Not implementation yet
func (ctx *Context) GetBlockHash(num uint64) (core.Word256, error) {
	return core.Zero256, nil
}
