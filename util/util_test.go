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

package util

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLog256(t *testing.T) {
	require.Equal(t, 0, Log256(big.NewInt(1)))
	require.Equal(t, 0, Log256(big.NewInt(255)))
	require.Equal(t, 1, Log256(big.NewInt(256)))
	require.Equal(t, 1, Log256(big.NewInt(65535)))
	require.Equal(t, 2, Log256(big.NewInt(65536)))
}

func TestFixBytesLength(t *testing.T) {
	require.Len(t, FixBytesLength(nil, 32), 32)
}
