package evm

import (
	"evm/core"
	"fmt"
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

// String return string of log
// Note: This should be used only for testing.
// TODO: A better String of remove it.
func (l *Log) String() string {
	if len(l.Topics) == 0 {
		return fmt.Sprintf("Address is %x, data is %x", l.Address.Bytes(), l.Data)
	}
	var topic = "["
	for i := range l.Topics {
		if i != 0 {
			topic += ","
		}
		topic += fmt.Sprintf("%x", l.Topics[i].Bytes())
	}
	topic += "]"
	return fmt.Sprintf("Address is %x, topic is %s, data is %x", l.Address.Bytes(), topic, l.Data)
}
