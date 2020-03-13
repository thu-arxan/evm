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

package evm

import (
	"github.com/thu-arxan/evm/core"
	"github.com/thu-arxan/evm/util"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Note: You can find the sample from https://ethereum.stackexchange.com/questions/760/how-is-the-address-of-an-ethereum-contract-computed
func TestDefaultCreateAddress(t *testing.T) {
	address, err := core.HexToAddress("0x6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	require.NoError(t, err)
	newAddress := defaultCreateAddress(address, 0, func(bytes []byte) Address {
		return core.AddressFromBytes(bytes)
	})
	require.Equal(t, "cd234a471b72ba2f1ccf0a70fcaba648a5eecd8d", fmt.Sprintf("%x", newAddress.Bytes()))
	newAddress = defaultCreateAddress(address, 1, func(bytes []byte) Address {
		return core.AddressFromBytes(bytes)
	})
	require.Equal(t, "343c43a37d37dff08ae8c4a11544c718abb4fcf8", fmt.Sprintf("%x", newAddress.Bytes()))
	t.Log()
}

// Note: You can find the sample from https://learnblockchain.cn/docs/eips/eip-1014.html#%E7%A4%BA%E4%BE%8B
func TestDefaultCreate2Address(t *testing.T) {
	address, err := core.HexToAddress("00000000000000000000000000000000deadbeef")
	require.NoError(t, err)
	salt, _ := util.HexToBytes("00000000000000000000000000000000000000000000000000000000cafebabe")
	code, _ := util.HexToBytes("deadbeef")
	newAddress := defaultCreate2Address(address, salt, code, func(bytes []byte) Address {
		return core.AddressFromBytes(bytes)
	})
	require.Equal(t, "60f3f640a8508fc6a86d45df051962668e1e8ac7", fmt.Sprintf("%x", newAddress.Bytes()))
}
