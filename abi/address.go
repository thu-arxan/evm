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

package abi

// This file provide a function which allow user define the address unpack they need.

var (
	addressLength   = 20
	addressToString func([]byte) string
	stringToAddress func(string) ([]byte, error)
)

// SetAddressParser set address reflect type and its length and toString function
func SetAddressParser(length int, toString func([]byte) string, toAddress func(string) ([]byte, error)) {
	addressLength = length
	addressToString = toString
	stringToAddress = toAddress
}

// ResetAddressParser reset reflect type and its length and toString function
func ResetAddressParser() {
	addressLength = 20
	addressToString = nil
	stringToAddress = nil
}
