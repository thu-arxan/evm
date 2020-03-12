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
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddressFromBytes(t *testing.T) {
	var bytes = randBytes(20)
	address := AddressFromBytes(bytes)
	require.Equal(t, bytes, address.Bytes())
	bytes = randBytes(12)
	address = AddressFromBytes(bytes)
	require.True(t, strings.HasSuffix(fmt.Sprintf("%x", address.Bytes()), fmt.Sprintf("%x", bytes)))
	bytes = randBytes(32)
	address = AddressFromBytes(bytes)
	require.True(t, strings.HasSuffix(fmt.Sprintf("%x", bytes), fmt.Sprintf("%x", address.Bytes())))
}
