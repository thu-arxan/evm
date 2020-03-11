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

package core

import (
	"evm/util"
	"strings"
)

// Address is Word160
type Address Word160

// Here defines some length
const (
	AddressLength    = Word160Length
	AddressHexLength = 2 * AddressLength
)

// ZeroAddress is zero
var ZeroAddress = Address{}

// Word256 return Word256 of address
func (address Address) Word256() Word256 {
	return Word160(address).Word256()
}

// Bytes return bytes of address, which length is 20
func (address Address) Bytes() []byte {
	return Word160(address).Bytes()
}

// AddressFromBytes returns an addres. It will cut left if len(bs) > 20, else add zeros at left if len(bs) < 20
func AddressFromBytes(bs []byte) (address Address) {
	if len(bs) > Word160Length {
		bs = bs[len(bs)-Word160Length:]
	} else if len(bs) < Word160Length {
		bs = util.BytesCombine(make([]byte, Word160Length-len(bs)), bs)
	}
	copy(address[:], bs)
	return
}

// HexToAddress convert hex string to Address, hex string could begin with 0x or 0X or nothing
func HexToAddress(s string) (address Address, err error) {
	if strings.HasPrefix(s, "0x") {
		s = strings.Replace(s, "0x", "", 1)
	} else if strings.HasPrefix(s, "0X") {
		s = strings.Replace(s, "0X", "", 1)
	}
	bytes, err := util.HexToBytes(s)
	if err != nil {
		return ZeroAddress, err
	}
	return AddressFromBytes(bytes), nil
}
