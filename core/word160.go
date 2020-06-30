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

// Here defines some consts
const (
	Word160Length       = 20
	Word256Word160Delta = 12
)

// Zero160 is the zero of Word160
var Zero160 = Word160{}

// Word160 is bytes which length is Word160Length
type Word160 [Word160Length]byte

// Word256 convert Word160 to Word256
// The function will add zeros before Word160 until its length == 32
func (w Word160) Word256() (word256 Word256) {
	copy(word256[Word256Word160Delta:], w[:])
	return
}

// Bytes return bytes of Word160
func (w Word160) Bytes() []byte {
	return w[:]
}
