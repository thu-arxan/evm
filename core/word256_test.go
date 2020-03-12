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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWord256Copy(t *testing.T) {
	var word256 Word256
	word256[0] = byte(1)
	copy := word256.Copy()
	require.Equal(t, word256.Bytes(), copy.Bytes())
	word256[0] = byte(2)
	require.NotEqual(t, word256.Bytes(), copy.Bytes())
}

func TestWord256ToWord160(t *testing.T) {
	var word160 Word160
	for i := range word160 {
		word160[i] = byte(util.RandNum(10))
	}
	word256 := word160.Word256()
	require.Equal(t, word160.Bytes(), word256.Word160().Bytes())
}
