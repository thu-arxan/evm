// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
package precompile

import (
	"errors"
)

// Contract is the basic interface for native Go contracts. The implementation
// requires a deterministic gas count based on the input size of the Run method of the
// contract.
type Contract interface {
	RequiredGas(input []byte) uint64  // RequiredPrice calculates the contract gas use
	Run(input []byte) ([]byte, error) // Run runs the precompiled contract
}

// IsPrecompile return if an address is precompile contract
func IsPrecompile(address []byte) bool {
	if len(address) == 0 {
		return false
	}
	for i := 0; i < len(address)-1; i++ {
		if address[i] != 0 {
			return false
		}
	}
	b := address[len(address)-1]
	if b >= 1 && b <= 9 {
		return true
	}
	return false
}

// New is the constructor of precompile contract
func New(address []byte) (Contract, error) {
	if !IsPrecompile(address) {
		return nil, errors.New("Not a precompile contract")
	}
	switch address[len(address)-1] {
	case 1:
		return &ecrecover{}, nil
	case 2:
		return &sha256hash{}, nil
	case 3:
		return &ripemd160hash{}, nil
	case 4:
		return &dataCopy{}, nil
	case 5:
		return &bigModExp{}, nil
	case 6:
		return &bn256AddIstanbul{}, nil
	case 7:
		return &bn256ScalarMulIstanbul{}, nil
	case 8:
		return &bn256PairingIstanbul{}, nil
	case 9:
		return &blake2F{}, nil
	default:
		return nil, errors.New("Not implementation yet")
	}
}
