//  Copyright 2020 The THU-Arxan Authors
//  This file is part of the evm library.
//
//  The evm library is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Lesser General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  The evm library is distributed in the hope that it will be useful,/
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
//  GNU Lesser General Public License for more details.
//
//  You should have received a copy of the GNU Lesser General Public License
//  along with the evm library. If not, see <http://www.gnu.org/licenses/>.
//

package errors

type code struct {
	err string
}

func newCode(err string) code {
	return code{
		err: err,
	}
}

func (c code) Error() string {
	return c.err
}

// Here defines some kind of error code
var (
	None                   = newCode("none")
	UnknownAddress         = newCode("unknown address")
	InsufficientBalance    = newCode("insufficient balance")
	InvalidJumpDest        = newCode("invalid jump destination")
	InsufficientGas        = newCode("insufficient gas")
	MemoryOutOfBounds      = newCode("memory out of bounds")
	CodeOutOfBounds        = newCode("code out of bounds")
	InputOutOfBounds       = newCode("input out of bounds")
	ReturnDataOutOfBounds  = newCode("data out of bounds")
	CallStackOverflow      = newCode("call stack overflow")
	CallStackUnderflow     = newCode("call stack underflow")
	DataStackOverflow      = newCode("data stack overflow")
	DataStackUnderflow     = newCode("data stack underflow")
	InvalidContract        = newCode("invalid contract")
	PermissionDenied       = newCode("permission denied")
	NativeContractCodeCopy = newCode("tried to copy native contract code")
	ExecutionAborted       = newCode("execution aborted")
	ExecutionReverted      = newCode("execution reverted")
	NativeFunction         = newCode("native function error")
	EventPublish           = newCode("event publish error")
	InvalidString          = newCode("invalid string")
	EventMapping           = newCode("event mapping error")
	Generic                = newCode("generic error")
	InvalidAddress         = newCode("invalid address")
	DuplicateAddress       = newCode("duplicate address")
	InsufficientFunds      = newCode("insufficient funds")
	Overpayment            = newCode("overpayment")
	ZeroPayment            = newCode("zero payment error")
	InvalidSequence        = newCode("invalid sequence number")
	ReservedAddress        = newCode("address is reserved for SNative or internal use")
	IllegalWrite           = newCode("callee attempted to illegally modify state")
	IntegerOverflow        = newCode("integer overflow")
	InvalidProposal        = newCode("proposal is invalid")
	ExpiredProposal        = newCode("proposal is expired since sequence number does not match")
	ProposalExecuted       = newCode("proposal has already been executed")
	NoInputPermission      = newCode("account has no input permission")
	InvalidBlockNumber     = newCode("invalid block number")
	BlockNumberOutOfRange  = newCode("block number out of range")
	AlreadyVoted           = newCode("vote already registered for this address")
	UnresolvedSymbols      = newCode("code has unresolved symbols")
	InvalidContractCode    = newCode("contract being created with unexpected code")
	NonExistentAccount     = newCode("account does not exist")
	UnknownOpcode          = newCode("unknown opcode")
)
