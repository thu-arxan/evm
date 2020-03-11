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
	"evm"
	"evm/db"
	"evm/example"
	"evm/util"
	"testing"

	"github.com/stretchr/testify/require"

)

var (
	createBin = "sols/Create_sol_C.bin"
	createAbi = "sols/Create_sol_C.abi"
	CCode  []byte
	CAddress evm.Address
)

func TestCreateSol(t *testing.T) {
	binBytes, err := util.ReadBinFile(createBin)
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var origin = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	var exceptAddress = `cd234a471b72ba2f1ccf0a70fcaba648a5eecd8d`
	CCode, CAddress = deployContract(t, memoryDB, bc, origin, binBytes, exceptAddress, "", 344952)
	callCreate(t, memoryDB, bc, origin, mustPack(createAbi, "createAndGetBalance", "44", "0"), 95396)
}

func callCreate(t *testing.T, db evm.DB, bc evm.Blockchain, caller evm.Address, payload []byte, gasCost uint64) {
	var gasQuota uint64 = 1000000
	var gas = gasQuota
	output, err := evm.New(bc, db, &evm.Context{
		Input: payload,
		Value: 0,
		Gas: &gas,
	}).Call(caller, CAddress, CCode)
	require.NoError(t, err)
	if gasCost != 0 {
		require.EqualValues(t, gasCost, gasQuota - gas)
	}

	t.Log(output)
	t.Log(gasQuota - gas)
}
