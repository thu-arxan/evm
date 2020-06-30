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
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWord160ToWord256(t *testing.T) {
	var word160 Word160
	for i := range word160 {
		word160[i] = byte(util.RandNum(10))
	}
	word160Hex := fmt.Sprintf("%x", word160.Bytes())
	word256Hex := fmt.Sprintf("%x", word160.Word256().Bytes())
	require.Equal(t, "000000000000000000000000"+word160Hex, word256Hex)
}
