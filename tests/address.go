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

package tests

import (
	"evm/util"
	"fmt"
)

// AddressLength define the length of Address
const AddressLength = 12

// Address is the address
type Address struct {
	value [AddressLength]byte
}

// HexToAddress convert hex string to address
func HexToAddress(s string) (*Address, error) {
	bytes, err := util.HexToBytes(s)
	if err != nil {
		return nil, err
	}
	bytes = util.FixBytesLength(bytes, AddressLength)
	var a Address
	for i := range a.value {
		a.value[i] = bytes[i]
	}
	return &a, nil
}

// BytesToAddress convert bytes to address
func BytesToAddress(bytes []byte) *Address {
	bytes = util.FixBytesLength(bytes, AddressLength)
	var a Address
	for i := range a.value {
		a.value[i] = bytes[i]
	}
	return &a
}

// Bytes return bytes of Address
func (a *Address) Bytes() []byte {
	var bytes = make([]byte, AddressLength)
	for i := range bytes {
		bytes[i] = a.value[i]
	}
	return bytes
}

// Length return the length of address
func (a *Address) Length() int {
	return AddressLength
}

// RandomAddress random generate an address
func RandomAddress() *Address {
	var a Address
	for i := range a.value {
		a.value[i] = byte(util.RandNum(128))
	}
	return &a
}

// String return hex of address
func (a *Address) String() string {
	return fmt.Sprintf("%X", a.value)
}
