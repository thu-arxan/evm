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

import "fmt"

// Pack provide a easy way to pack, it is a simple wrapper of New & PackValues
func Pack(abiFile, funcName string, inputs ...string) ([]byte, error) {
	abi, err := New(abiFile)
	if err != nil {
		return nil, err
	}
	return abi.PackValues(funcName, inputs...)
}

// Unpack provide a easy way to unpack, it is a simple wrapper of New & UnpackValues
func Unpack(abiFile, funcName string, data []byte) (values []string, err error) {
	defer func() {
		if e := recover(); e != nil {
			values = nil
			err = fmt.Errorf("unpack failed because %v", e)
		}
	}()
	abi, err := New(abiFile)
	if err != nil {
		return nil, err
	}
	values, err = abi.UnpackValues(funcName, data)
	return
}
