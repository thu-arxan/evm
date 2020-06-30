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

package example

import "evm/util"

// Address is the address
type Address [20]byte

// Bytes is the implementation of interface
func (a *Address) Bytes() []byte {
	return a[:]
}

// Copy return the copy of address
func (a *Address) Copy() *Address {
	var ret Address
	copy(ret[:], a.Bytes())
	return &ret
}

// BytesToAddress convert bytes to address
func BytesToAddress(bytes []byte) *Address {
	var a Address
	copy(a[:], util.FixBytesLength(bytes, 20))
	return &a
}

// HexToAddress convert hex string to address, string may begin with 0x, 0X or nothing
func HexToAddress(hex string) *Address {
	var a Address
	if bytes, err := util.HexToBytes(hex); err == nil {
		copy(a[:], util.FixBytesLength(bytes, 20))
	}
	return &a
}

// RandomAddress return a random address
func RandomAddress() *Address {
	var a Address
	for i := range a {
		a[i] = byte(util.RandNum(128))
	}
	return &a
}

// ZeroAddress return zero address
func ZeroAddress() *Address {
	var a Address
	return &a
}
