package evm

import (
	"evm/core"
)

// Log is the log of evm
type Log struct {
	// Consensus field
	Address Address `json:"address"`
	// list of topics provided by the contract.
	Topics []core.Word256 `json:"topics"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data"`

	// Derived fields, so the database should record the context to support these fields
	BlockNumber uint64 `json:"blockNumber"`
	// hash of the transaction
	TxHash []byte `json:"transactionHash"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex"`
	// hash of the block in which the transaction was included
	BlockHash []byte `json:"blockHash"`
	// index of the log in the block
	Index uint `json:"logIndex"`
}
