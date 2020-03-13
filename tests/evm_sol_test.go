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
	"github.com/thu-arxan/evm"
	"github.com/thu-arxan/evm/db"
	"github.com/thu-arxan/evm/example"
	"github.com/thu-arxan/evm/util"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	evmCodeBin = "sols/Ethereum_sol_OpCodes.bin"
	evmCodeAbi = "sols/Ethereum_sol_OpCodes.abi"
	evmCode []byte
	evmCodeAddress evm.Address
)

func TestEvm(t *testing.T) {
	binBytes, err := util.ReadBinFile(evmCodeBin)
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var origin = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	evmCode, evmCodeAddress = deployContract(t, memoryDB, bc, origin, binBytes, "", "", 388049)
	input := mustPack(evmCodeAbi, "test")
	var gas uint64 = 1000000
	output, err := evm.New(bc, memoryDB, &evm.Context{
		Input: input,
		Value: 0,
		Gas: &gas,
		BlockHeight: 1,
	}).Call(origin, evmCodeAddress, evmCode)
	require.NoError(t, err)
	t.Log(output)
	t.Log(gas)
}