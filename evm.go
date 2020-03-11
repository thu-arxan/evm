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
	"bytes"
	"evm/util/math"
	"math/big"

	"evm/core"
	"evm/errors"
	"evm/gas"
	"evm/precompile"
	"evm/util"

	"evm/crypto"

	"github.com/sirupsen/logrus"
)

var (
	log   = logrus.WithFields(logrus.Fields{"package": "evm"})
	debug = false
)

var (
	tt255 = math.BigPow(2, 255)
)

// Here defines some default stack capacity variables
const (
	DefaultStackCapacity    uint64 = 1024
	DefaultMaxStackCapacity uint64 = 32 * 1024
	MaxCodeSize             int    = 24576
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

// SetDebug set debug and logrus log level
func SetDebug(isDebug bool) {
	debug = isDebug
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

// EVM is the evm
type EVM struct {
	origin         Address
	ctx            *Context
	bc             Blockchain
	cache          *Cache
	memoryProvider func(errorSink errors.Sink) Memory
	stackDepth     uint64
	refund         uint64
	sync           bool
}

// New is the constructor of EVM
func New(bc Blockchain, db DB, ctx *Context) *EVM {
	return &EVM{
		bc:             bc,
		cache:          NewCache(db),
		memoryProvider: DefaultDynamicMemoryProvider,
		ctx:            ctx,
		sync:           true,
	}
}

// Create create a contract account, and return an error if there exist a contract on the address
func (evm *EVM) Create(caller Address) ([]byte, Address, error) {
	if evm.origin == nil {
		evm.origin = caller
	}
	if len(evm.ctx.Input) == 0 {
		return nil, nil, errors.InvalidContractCode
	}
	nonce := evm.cache.GetNonce(caller)
	address := evm.bc.CreateAddress(caller, nonce)
	// call default implementaion if the user do no want to implement it
	if address == nil {
		address = defaultCreateAddress(caller, evm.cache.GetNonce(caller), evm.bc.BytesToAddress)
	}
	if err := evm.createAccount(caller, address); err != nil {
		return nil, address, err
	}
	// update caller nonce and update
	callerAccount := evm.cache.GetAccount(caller)
	callerAccount.SetNonce(nonce + 1)
	evm.cache.UpdateAccount(callerAccount)
	// set contract nonce -> 1 and update
	contract := evm.cache.GetAccount(address)
	contract.SetNonce(1)
	evm.cache.UpdateAccount(contract)
	// transfer and run
	if err := evm.transfer(caller, address, evm.ctx.Value); err != nil {
		return nil, nil, err
	}
	code, err := evm.callWithDepth(caller, address, evm.ctx.Input)
	if err != nil {
		return nil, nil, err
	}
	createDataGas := uint64(len(code)) * gas.CreateData
	if useGasNegative(evm.ctx.Gas, createDataGas) != nil {
		return nil, nil, errors.InsufficientGas
	}
	if len(code) > MaxCodeSize {
		return nil, nil, errors.CodeOutOfBounds
	}
	contract = evm.cache.GetAccount(address)
	contract.SetCode(code)

	if err := evm.cache.UpdateAccount(contract); err != nil {
		return nil, nil, err
	}

	if evm.sync {
		evm.cache.Sync()
	}
	return code, address, nil
}

// Call run code on evm, and it will sync change to db if error is nil
func (evm *EVM) Call(caller, callee Address, code []byte) ([]byte, error) {
	if evm.origin == nil {
		evm.origin = caller
	}
	if err := evm.transfer(caller, callee, evm.ctx.Value); err != nil {
		return nil, err
	}

	return evm.CallWithoutTransfer(caller, callee, code)
}

// CallWithoutTransfer is call without transfer, and it will sync change to db if error is nil
func (evm *EVM) CallWithoutTransfer(caller, callee Address, code []byte) (output []byte, err error) {
	if evm.origin == nil {
		evm.origin = caller
	}
	if precompile.IsPrecompile(callee.Bytes()) {
		contract, err := precompile.New(callee.Bytes())
		if err != nil {
			return nil, err
		}
		if err := useGasNegative(evm.ctx.Gas, contract.RequiredGas(evm.ctx.Input)); err != nil {
			return nil, err
		}
		output, err = contract.Run(evm.ctx.Input)
	} else {
		output, err = evm.callWithDepth(caller, callee, code)
		if err != nil {
			return
		}
	}

	// sync change to db if no error
	if evm.sync {
		evm.cache.Sync()
	}
	return
}

// GetRefund return the refund
func (evm *EVM) GetRefund() uint64 {
	return evm.refund
}

func (evm *EVM) addRefund(gas uint64) {
	evm.refund += gas
}

func (evm *EVM) subRefund(gas uint64) {
	evm.refund -= gas
}

func (evm *EVM) transfer(caller, callee Address, value uint64) error {
	if value == 0 {
		return nil
	}

	from := evm.cache.GetAccount(caller)
	if err := from.SubBalance(value); err != nil {
		return err
	}

	to := evm.cache.GetAccount(callee)
	if err := to.AddBalance(value); err != nil {
		return err
	}

	if err := evm.cache.UpdateAccount(from); err != nil {
		return err
	}

	if err := evm.cache.UpdateAccount(to); err != nil {
		return err
	}

	return nil
}

func (evm *EVM) callWithDepth(caller, callee Address, code []byte) ([]byte, error) {
	if len(code) > 0 {
		evm.stackDepth++
		if evm.stackDepth > 1024 {
			return nil, errors.CallStackOverflow
		}
		output, err := evm.call(caller, callee, code)
		evm.stackDepth--
		return output, err
	}
	return nil, nil
}

// call does not transfer 'value' or modify the callDepth.
func (evm *EVM) call(caller, callee Address, code []byte) ([]byte, error) {
	var maybe = errors.NewMaybe()
	var ctx = evm.ctx
	var pc uint64
	var stack = NewStack(DefaultStackCapacity, DefaultMaxStackCapacity, ctx.Gas, maybe, evm.bc.BytesToAddress)
	var memory = evm.memoryProvider(maybe)

	var returnData []byte

	for {
		if maybe.Error() != nil {
			if maybe.Error() != errors.ExecutionReverted {
				*ctx.Gas = 0
			}
			return nil, maybe.Error()
		}

		var op = getOpCode(code, pc)
		if debug {
			log.Debugf("(pc) %-3d (op) %-14s (st) %-4d (gas) %d", pc, op.String(), stack.Len(), *ctx.Gas)
		}

		switch op {
		case ADD: // 0x01
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopBigInt(), stack.PeekBigInt()
			if debug {
				log.Debugf("  %v + %v ", x, y)
			}
			math.U256(y.Add(x, y))

		case MUL: // 0x02
			maybe.PushError(useGasNegative(ctx.Gas, gas.Low))
			x, y := stack.PopBigInt(), stack.PeekBigInt()
			if debug {
				log.Debugf("  %v * %v ", x, y)
			}
			math.U256(y.Mul(x, y))

		case SUB: // 0x03
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopBigInt(), stack.PeekBigInt()
			if debug {
				log.Debugf("  %v - %v ", x, y)
			}
			math.U256(y.Sub(x, y))

		case DIV: // 0x04
			maybe.PushError(useGasNegative(ctx.Gas, gas.Low))
			x, y := stack.PopBigInt(), stack.PeekBigInt()
			if debug {
				log.Debugf("  %v / %v ", op, y)
			}
			if y.Sign() != 0 {
				math.U256(y.Div(x, y))
			} else {
				y.SetUint64(0)
			}

		case SDIV: // 0x05
			maybe.PushError(useGasNegative(ctx.Gas, gas.Low))
			x, y := math.S256(stack.PopBigInt()), math.S256(stack.PeekBigInt())
			if debug {
				log.Debugf("  %v / %v ", x, y)
			}
			if y.Sign() == 0 || x.Sign() == 0 {
				y.SetUint64(0)
			} else {
				if x.Sign() != y.Sign() {
					y.Div(x.Abs(x), y.Abs(y))
					y.Neg(y)
				} else {
					y.Div(x.Abs(x), y.Abs(y))
				}
				math.U256(y)
			}

		case MOD: // 0x06
			maybe.PushError(useGasNegative(ctx.Gas, gas.Low))
			x, y := stack.PopBigInt(), stack.PeekBigInt()
			if debug {
				log.Debugf("  %v %% %v", x, y)
			}
			if y.Sign() == 0 {
				y.SetUint64(0)
			} else {
				math.U256(y.Mod(x, y))
			}

		case SMOD: // 0x07
			maybe.PushError(useGasNegative(ctx.Gas, gas.Low))
			x, y := math.S256(stack.PopBigInt()), math.S256(stack.PeekBigInt())
			if debug {
				log.Debugf("  %v %% %v", x, y)
			}
			if y.Sign() == 0 {
				y.SetUint64(0)
			} else {
				if x.Sign() < 0 {
					y.Mod(x.Abs(x), y.Abs(y))
					y.Neg(y)
				} else {
					y.Mod(x.Abs(x), y.Abs(y))
				}
				math.U256(y)
			}

		case ADDMOD: // 0x08
			maybe.PushError(useGasNegative(ctx.Gas, gas.Mid))
			x, y, z := stack.PopBigInt(), stack.PopBigInt(), stack.PeekBigInt()
			if debug {
				log.Debugf("  %v %% %v", x, y)
			}
			if z.Sign() == 0 {
				z.SetUint64(0)
			} else {
				x.Add(x, y)
				z.Mod(x, z)
				math.U256(z)
			}

		case MULMOD: // 0x09
			maybe.PushError(useGasNegative(ctx.Gas, gas.Mid))
			x, y, z := stack.PopBigInt(), stack.PopBigInt(), stack.PeekBigInt()
			if debug {
				log.Debugf("  %v %% %v", x, y)
			}
			if z.Sign() == 0 {
				// stack.PushBigInt(x.SetUint64(0))
				z.SetUint64(0)
			} else {
				x.Mul(x, y)
				z.Mod(x, z)
				math.U256(z)
			}

		case EXP: // 0x0A
			base, exponent := stack.PopBigInt(), stack.PeekBigInt()
			if debug {
				log.Debugf("  %v ** %v", base, exponent)
			}
			if exponent.Sign() == 0 {
				maybe.PushError(useGasNegative(ctx.Gas, gas.Exp))
			} else {
				maybe.PushError(useGasNegative(ctx.Gas, gas.Exp+gas.ExpByte*uint64(1+util.Log256(exponent))))
			}
			cmpToOne := exponent.Cmp(math.Big1)
			if cmpToOne < 0 { // Exponent is zero
				// x ^ 0 == 1
				exponent.SetUint64(1)
			} else if base.Sign() == 0 {
				// 0 ^ y, if y != 0, == 0
				exponent.SetUint64(0)
			} else if cmpToOne == 0 { // Exponent is one
				// x ^ 1 == x
			} else {
				exponent = math.Exp(base, exponent)
			}

		case SIGNEXTEND: // 0x0B
			maybe.PushError(useGasNegative(ctx.Gas, gas.Low))
			back := stack.PopBigInt()
			if back.Cmp(big.NewInt(31)) < 0 {
				bit := uint(back.Uint64()*8 + 7)
				res := stack.PopBigInt()
				num := res
				mask := back.Lsh(core.Big1, bit)
				mask.Sub(mask, core.Big1)
				if res.Bit(int(bit)) > 0 {
					res.Or(res, mask.Not(mask))
				} else {
					res.And(res, mask)
				}
				stack.PushBigInt(res)
				if debug {
					log.Debugf("  %v signextend %v = %v", num, back, res)
				}
			}

		case LT: // 0x10
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopBigInt(), stack.PeekBigInt()
			if debug {
				log.Debugf("  %v < %v", x, y)
			}
			if x.Cmp(y) < 0 {
				y.SetUint64(1)
			} else {
				y.SetUint64(0)
			}

		case GT: // 0x11
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopBigInt(), stack.PeekBigInt()
			if debug {
				log.Debugf("  %v > %v", x, y)
			}
			if x.Cmp(y) > 0 {
				y.SetUint64(1)
			} else {
				y.SetUint64(0)
			}

		case SLT: // 0x12
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopBigInt(), stack.PeekBigInt()

			xSign := x.Cmp(tt255)
			ySign := y.Cmp(tt255)
			if debug {
				log.Debugf("  %v < %v ", x, y)
			}
			switch {
			case xSign >= 0 && ySign < 0:
				y.SetUint64(1)
			case xSign < 0 && ySign >= 0:
				y.SetUint64(0)
			default:
				if x.Cmp(y) < 0 {
					y.SetUint64(1)
				} else {
					y.SetUint64(0)
				}
			}

		case SGT: // 0x13
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopBigInt(), stack.PeekBigInt()
			if debug {
				log.Debugf("  %v > %v", x, y)
			}
			xSign := x.Cmp(tt255)
			ySign := y.Cmp(tt255)

			switch {
			case xSign >= 0 && ySign < 0:
				y.SetUint64(0)
			case xSign < 0 && ySign >= 0:
				y.SetUint64(1)
			default:
				if x.Cmp(y) > 0 {
					y.SetUint64(1)
				} else {
					y.SetUint64(0)
				}
			}

		case EQ: // 0x14
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopBigInt(), stack.PeekBigInt()
			if debug {
				log.Debugf("  %v == %v", x, y)
			}
			if x.Cmp(y) == 0 {
				y.SetUint64(1)
			} else {
				y.SetUint64(0)
			}

		case ISZERO: // 0x15
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x := stack.PeekBigInt()
			if debug {
				log.Debugf("  %v == 0", x)
			}
			if x.Sign() > 0 {
				x.SetUint64(0)
			} else {
				x.SetUint64(1)
			}

		case AND: // 0x16
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopBigInt(), stack.PeekBigInt()
			y.And(x, y)
			if debug {
				log.Debugf("  %v & %v", x, y)
			}

		case OR: // 0x17
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopBigInt(), stack.PeekBigInt()
			y.Or(x, y)
			if debug {
				log.Debugf("  %v | %v", x, y)
			}

		case XOR: // 0x18
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopBigInt(), stack.PeekBigInt()
			y.Xor(x, y)
			if debug {
				log.Debugf("  %v ^ %v", x, y)
			}

		case NOT: // 0x19
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x := stack.PeekBigInt()
			math.U256(x.Not(x))
			if debug {
				log.Debugf("  !%v", x)
			}

		case BYTE: // 0x1A
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			th, val := stack.PopBigInt(), stack.PeekBigInt()
			if th.Cmp(math.Big32) < 0 {
				b := math.Byte(val, 32, int(th.Int64()))
				val.SetUint64(uint64(b))
			} else {
				val.SetUint64(0)
			}
			if debug {
				log.Debugf("  0x%X", val.Bytes())
			}

		case SHL: //0x1B
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			shift, value := math.U256(stack.PopBigInt()), math.U256(stack.PeekBigInt())
			if shift.Cmp(math.Big256) >= 0 {
				value.SetUint64(0)
			} else {
				n := uint(shift.Uint64())
				math.U256(value.Lsh(value, n))
			}
			if debug {
				log.Debugf("  %v << %v", value, shift)
			}

		case SHR: //0x1C
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			shift, x := stack.PopBigInt(), stack.PeekBigInt()
			if debug {
				log.Debugf("  %v << %v", x, shift)
			}
			if shift.Cmp(math.Big256) >= 0 {
				x.SetUint64(0)
			} else {
				x.Rsh(x, uint(shift.Uint64()))
			}

		case SAR: //0x1D
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			shift, x := stack.PopBigInt(), math.S256(stack.PeekBigInt())
			if debug {
				log.Debugf("  %v << %v", x, shift)
			}
			if shift.Cmp(core.Big256) >= 0 {
				if x.Sign() < 0 {
					x.SetInt64(-1)
				} else {
					x.SetInt64(0)
				}
			} else {
				x.Rsh(x, uint(shift.Uint64()))
			}

		case SHA3: // 0x20
			offset, size := stack.PopBigInt(), stack.PopBigInt()
			maybe.PushError(useGasNegative(ctx.Gas, gas.SHA3+gas.SHA3Word*((size.Uint64()+31)/32)))
			data, memoryGas := memory.Read(offset, size)
			maybe.PushError(useGasNegative(ctx.Gas, memoryGas))
			data = crypto.Keccak256(data)
			stack.PushBytes(data)
			if debug {
				log.Debugf(" (%v) %X", size, data)
			}

		case ADDRESS: // 0x30
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushAddress(callee)
			if debug {
				log.Debugf("  %v", callee)
			}

		case BALANCE: // 0x31
			// todo: we may peek and set uint?
			maybe.PushError(useGasNegative(ctx.Gas, gas.Balance))
			address := stack.PopAddress()
			balance := evm.getAccount(address).GetBalance()
			stack.PushUint64(balance)
			if debug {
				log.Debugf("  %v (%v)", balance, address)
			}

		case ORIGIN: // 0x32
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushAddress(evm.origin)
			if debug {
				log.Debugf("  %v", evm.origin)
			}

		case CALLER: // 0x33
			// todo:2x time
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushAddress(caller)
			if debug {
				log.Debugf("  %v", caller)
			}

		case CALLVALUE: // 0x34
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushUint64(ctx.Value)
			if debug {
				log.Debugf("  %v", ctx.Value)
			}

		case CALLDATALOAD: // 0x35
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			offset := stack.PeekBigInt()
			data, err := util.SubSlice(ctx.Input, offset.Uint64(), 32)
			if err != nil {
				maybe.PushError(errors.InputOutOfBounds)
			}
			res := core.LeftPadWord256(data)
			offset.SetBytes(res.Bytes())
			if debug {
				log.Debugf("  0x%v", res)
			}

		case CALLDATASIZE: // 0x36
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushUint64(uint64(len(ctx.Input)))
			if debug {
				log.Debugf("  %d", len(ctx.Input))
			}

		case CALLDATACOPY: // 0x37
			memOff := stack.PopBigInt()
			inputOff := stack.PopBigInt()
			length := stack.PopBigInt()
			data := util.GetDataBig(ctx.Input, inputOff, length)
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			gasCost := memory.Write(memOff, data) + wordGas(length.Uint64(), gas.Copy)
			maybe.PushError(useGasNegative(ctx.Gas, gasCost))
			if debug {
				log.Debugf("  [%v, %v, %v] %X", memOff, inputOff, length, data)
			}

		case CODESIZE: // 0x38
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			l := uint64(len(code))
			stack.PushUint64(l)
			if debug {
				log.Debugf("  %d", l)
			}

		case CODECOPY: // 0x39
			memOff := stack.PopBigInt()
			codeOff := stack.PopBigInt()
			length := stack.PopBigInt()
			data := util.GetDataBig(code, codeOff, length)
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			gasCost := memory.Write(memOff, data) + wordGas(length.Uint64(), gas.Copy)
			maybe.PushError(useGasNegative(ctx.Gas, gasCost))
			if debug {
				log.Debugf("  [%v, %v, %v] %X", memOff, codeOff, length, data)
			}

		case GASPRICE: // 0x3A
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushUint64(ctx.GasPrice)
			if debug {
				log.Debugf("  %v", ctx.GasPrice)
			}

		case EXTCODESIZE: // 0x3B
			maybe.PushError(useGasNegative(ctx.Gas, gas.ExtCode))
			address := stack.PopAddress()
			acc := evm.getAccount(address)
			length := uint64(len(acc.GetCode()))
			stack.PushUint64(length)
			if debug {
				log.Debugf("  %d", length)
			}

		case EXTCODECOPY: // 0x3C
			maybe.PushError(useGasNegative(ctx.Gas, gas.ExtCode))
			address := stack.PopAddress()
			code := evm.getAccount(address).GetCode()
			memOff := stack.PopBigInt()
			codeOff := stack.PopBigInt()
			length := stack.PopBigInt()
			data := util.GetDataBig(code, codeOff, length)
			gasCost := memory.Write(memOff, data) + wordGas(length.Uint64(), gas.Copy)
			maybe.PushError(useGasNegative(ctx.Gas, gasCost))
			if debug {
				log.Debugf("  [%v, %v, %v] %X", memOff, codeOff, length, data)
			}

		case RETURNDATASIZE: // 0x3D
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushUint64(uint64(len(returnData)))
			if debug {
				log.Debugf("  %d", len(returnData))
			}

		case RETURNDATACOPY: // 0x3E
			memOff, outputOff, length := stack.PopBigInt(), stack.PopBigInt(), stack.PopBigInt()
			end := new(big.Int).Add(outputOff, length)

			if !end.IsUint64() || uint64(len(returnData)) < end.Uint64() {
				maybe.PushError(errors.ReturnDataOutOfBounds)
				continue
			}
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow+gas.Copy*((length.Uint64()+31)/32)))
			gasCost := memory.Write(memOff, returnData[outputOff.Uint64():end.Uint64()]) + wordGas(length.Uint64(), gas.Copy)
			maybe.PushError(useGasNegative(ctx.Gas, gasCost))
			if debug {
				log.Debugf("  [%v, %v, %v] %X", memOff, outputOff, length, returnData)
			}

		case EXTCODEHASH: // 0x3F
			maybe.PushError(useGasNegative(ctx.Gas, gas.ExtcodeHash))
			address := stack.PopAddress()
			acc := evm.getAccount(address)
			// keccak256 hash of a contract's code
			var extcodehash core.Word256
			if isEmptyAccount(acc) {
				extcodehash = core.Zero256
			} else {
				if len(acc.GetCodeHash()) > 0 {
					copy(extcodehash[:], acc.GetCodeHash())
				} else {
					copy(extcodehash[:], crypto.Keccak256(acc.GetCode()))
				}
			}
			stack.Push(extcodehash)

		case BLOCKHASH: // 0x40
			maybe.PushError(useGasNegative(ctx.Gas, gas.BlockHash))
			blockNumber := stack.PopUint64()
			// Note: Here is >= other than > because block is not generated while running tx
			if blockNumber >= ctx.BlockHeight {
				if debug {
					log.Debugf("  attempted to get block hash of a non-existent block: %v", blockNumber)
				}
				stack.Push(core.Zero256)
			} else if ctx.BlockHeight-blockNumber > 256 {
				if debug {
					log.Debugf("  attempted to get block hash of a block %d outof range", blockNumber)
				}
				stack.Push(core.Zero256)
			} else {
				blockHash := evm.bc.GetBlockHash(blockNumber)
				stack.Push(core.LeftPadWord256(blockHash))
				if debug {
					log.Debugf("  0x%v", blockHash)
				}
			}

		case COINBASE: // 0x41
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushBytes(ctx.CoinBase)
			if debug {
				log.Debugf("  0x%v (NOT SUPPORTED)", stack.Peek())
			}

		case TIMESTAMP: // 0x42
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			blockTime := ctx.BlockTime
			stack.PushUint64(uint64(blockTime))
			if debug {
				log.Debugf("  %d", blockTime)
			}

		case NUMBER: // 0x43
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			number := ctx.BlockHeight
			stack.PushUint64(number)
			if debug {
				log.Debugf("  %d", number)
			}

		case DIFFICULTY: // Note: New version deprecated
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			difficulty := ctx.Difficulty
			stack.PushUint64(difficulty)
			if debug {
				log.Debugf("  %d", difficulty)
			}

		case GASLIMIT: // 0x45
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushUint64(ctx.GasLimit)
			if debug {
				log.Debugf("  %v", ctx.GasLimit)
			}

		case CHAINID: // 0x46
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.Push(core.Word256{})
			if debug {
				log.Debugf("  Not implemented")
			}

		case SELFBALANCE: // 0x47
			maybe.PushError(useGasNegative(ctx.Gas, gas.Low))
			balance := evm.getAccount(callee).GetBalance()
			stack.PushUint64(balance)
			if debug {
				log.Debugf("  %v (%v)", balance, callee)
			}

		case POP: // 0x50
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			popped := stack.PopBigInt()
			if debug {
				log.Debugf("  0x%v", popped.Bytes())
			}

		case MLOAD: // 0x51
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			offset := stack.PeekBigInt()
			data, memoryGas := memory.Read(offset, core.BigWord256Bytes)
			maybe.PushError(useGasNegative(ctx.Gas, memoryGas))
			// TODO: We do not need word256 as middle variable
			offset.SetBytes(core.LeftPadWord256(data).Bytes())
			if debug {
				log.Debugf("  0x%X @ 0x%v", data, offset)
			}

		case MSTORE: // 0x52
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			offset, val := stack.PopBigInt(), stack.PopBigInt()
			var data = make([]byte, 32)
			math.ReadBits(val, data)
			gasCost := memory.Write(offset, data)
			maybe.PushError(useGasNegative(ctx.Gas, gasCost))
			if debug {
				log.Debugf("  0x%v @ 0x%v", data, offset)
			}

		case MSTORE8: // 0x53
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			offset := stack.PopBigInt()
			val64 := stack.PopUint64()
			val := byte(val64 & 0xFF)
			gasCost := memory.Write(offset, []byte{val})
			maybe.PushError(useGasNegative(ctx.Gas, gasCost))
			if debug {
				log.Debugf("  [%v] 0x%X", offset, val)
			}

		case SLOAD: // 0x54
			// TODO: SLOAD is too slow!!!
			maybe.PushError(useGasNegative(ctx.Gas, gas.Sload))
			loc := stack.PeekBigInt()
			value := evm.cache.GetStorage(callee, core.BytesToWord256(loc.Bytes()))
			data := core.LeftPadWord256(value)
			loc.SetBytes(data.Bytes())
			if debug {
				log.Debugf("  %v {0x%v = 0x%v}", callee, loc, data)
			}

		case SSTORE: // 0x55
			loc, data := stack.Pop(), stack.Pop()
			currentData := evm.cache.GetStorage(callee, loc)
			if *ctx.Gas <= gas.SstoreSentryEIP2200 {
				maybe.PushError(errors.InsufficientGas)
			}
			if isEqual(data.Bytes(), currentData) {
				maybe.PushError(useGasNegative(ctx.Gas, gas.SstoreNoopEIP2200))
			} else {
				originData := evm.cache.db.GetStorage(callee, loc.Bytes())
				if isEqual(originData, currentData) {
					if isEmptyValue(originData) {
						maybe.PushError(useGasNegative(ctx.Gas, gas.SstoreInitEIP2200))
					} else {
						if isEmptyValue(data.Bytes()) {
							evm.addRefund(gas.SstoreClearRefundEIP2200)
						}
						maybe.PushError(useGasNegative(ctx.Gas, gas.SstoreCleanEIP2200))
					}
				} else {
					if !isEmptyValue(originData) {
						if isEmptyValue(currentData) { // recreate slot (2.2.1.1)
							evm.subRefund(gas.SstoreClearRefundEIP2200)
						} else if isEmptyValue(data.Bytes()) { // delete slot (2.2.1.2)
							evm.addRefund(gas.SstoreClearRefundEIP2200)
						}
					}
					if isEqual(originData, data.Bytes()) {
						if isEmptyValue(originData) { // reset to original inexistent slot (2.2.2.1)
							evm.addRefund(gas.SstoreInitRefundEIP2200)
						} else { // reset to original existing slot (2.2.2.2)
							evm.addRefund(gas.SstoreCleanRefundEIP2200)
						}
					}
					maybe.PushError(useGasNegative(ctx.Gas, gas.SstoreDirtyEIP2200))
				}
			}
			evm.cache.SetStorage(callee, loc, data.Bytes())
			if debug {
				log.Debugf("  %v {%v := %v}", callee, loc, data)
			}

		case JUMP: // 0x56
			maybe.PushError(useGasNegative(ctx.Gas, gas.Mid))
			to := stack.PopUint64()
			maybe.PushError(jump(code, to, &pc))
			continue

		case JUMPI: // 0x57
			maybe.PushError(useGasNegative(ctx.Gas, gas.High))
			pos := stack.PopUint64()
			cond := stack.Pop()
			if !cond.IsZero() {
				maybe.PushError(jump(code, pos, &pc))
				continue
			} else {
				if debug {
					log.Debugf(" false")
				}
			}

		case PC: // 0x58
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushUint64(pc)

		case MSIZE: // 0x59
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			// Note: Solidity will write to this offset expecting to find guaranteed
			// free memory to be allocated for it if a subsequent MSTORE is made to
			// this offset.
			capacity := memory.Capacity()
			stack.PushBigInt(capacity)
			if debug {
				log.Debugf("  0x%X", capacity)
			}

		case GAS: // 0x5A
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushUint64(*ctx.Gas)
			if debug {
				log.Debugf("  %X", *ctx.Gas)
			}

		case JUMPDEST: // 0x5B
			maybe.PushError(useGasNegative(ctx.Gas, gas.JumpDest))
			// Do nothing

		case PUSH1, PUSH2, PUSH3, PUSH4, PUSH5, PUSH6, PUSH7, PUSH8, PUSH9, PUSH10, PUSH11, PUSH12, PUSH13, PUSH14, PUSH15, PUSH16, PUSH17, PUSH18, PUSH19, PUSH20, PUSH21, PUSH22, PUSH23, PUSH24, PUSH25, PUSH26, PUSH27, PUSH28, PUSH29, PUSH30, PUSH31, PUSH32:
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			a := uint64(op - PUSH1 + 1)
			codeSegment, err := util.SubSlice(code, pc+1, a)
			if err != nil {
				maybe.PushError(errors.InputOutOfBounds)
			}
			res := core.LeftPadWord256(codeSegment)
			stack.Push(res)
			pc += a
			if debug {
				log.Debugf("  0x%v", res)
			}

		case DUP1, DUP2, DUP3, DUP4, DUP5, DUP6, DUP7, DUP8, DUP9, DUP10, DUP11, DUP12, DUP13, DUP14, DUP15, DUP16:
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			n := int(op - DUP1 + 1)
			stack.Dup(n)
			if debug {
				log.Debugf("  [%d] 0x%v", n, stack.Peek())
			}

		case SWAP1, SWAP2, SWAP3, SWAP4, SWAP5, SWAP6, SWAP7, SWAP8, SWAP9, SWAP10, SWAP11, SWAP12, SWAP13, SWAP14, SWAP15, SWAP16:
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			n := int(op - SWAP1 + 2)
			stack.Swap(n)
			if debug {
				log.Debugf("  [%d] %v", int(op-SWAP1+2), 0)
			}

		case LOG0, LOG1, LOG2, LOG3, LOG4:
			n := int(op - LOG0)
			topics := make([]core.Word256, n)
			offset, size := stack.PopBigInt(), stack.PopBigInt()
			for i := 0; i < n; i++ {
				topics[i] = stack.Pop()
			}
			data, memoryGas := memory.Read(offset, size)
			maybe.PushError(useGasNegative(ctx.Gas, memoryGas))
			maybe.PushError(useGasNegative(ctx.Gas, gas.Log+gas.LogData*size.Uint64()+uint64(op-LOG0)*gas.LogTopic))
			evm.cache.AddLog(&Log{
				Address: callee,
				Topics:  topics,
				Data:    data,
			})
			if debug {
				log.Debugf("  T:%v D:%X", topics, data)
			}

		case CREATE, CREATE2: // 0xF0, 0xFB
			if err := useGasNegative(ctx.Gas, gas.Create); err != nil {
				return nil, err
			}
			returnData = nil
			contractValue := stack.PopUint64()
			offset, size := stack.PopBigInt(), stack.PopBigInt()
			input, memoryGas := memory.Read(offset, size)
			if op == CREATE2 {
				memoryGas += wordGas(size.Uint64(), gas.SHA3Word)
			}
			// apply EIP150
			if err := useGasNegative(ctx.Gas, memoryGas); err != nil {
				return nil, err
			}
			gasPrev := *ctx.Gas / 64
			*ctx.Gas -= gasPrev

			var newAccountAddress Address
			if op == CREATE {
				newAccountAddress = evm.bc.CreateAddress(callee, evm.cache.GetNonce(callee))
				if newAccountAddress == nil {
					newAccountAddress = defaultCreateAddress(callee, evm.cache.GetNonce(callee), evm.bc.BytesToAddress)
				}
				calleeAccount := evm.cache.GetAccount(callee)
				calleeAccount.SetNonce(evm.cache.GetNonce(callee) + 1)
				maybe.PushError(evm.cache.UpdateAccount(calleeAccount))
			} else if op == CREATE2 {
				salt := stack.Pop()
				code := evm.getAccount(callee).GetCode()
				newAccountAddress = evm.bc.Create2Address(callee, salt.Bytes(), code)
				if newAccountAddress == nil {
					newAccountAddress = defaultCreate2Address(callee, salt.Bytes(), code, evm.bc.BytesToAddress)
				}
			}

			if evm.cache.Exist(newAccountAddress) {
				maybe.PushError(errors.InvalidAddress)
			}

			newAccount := evm.bc.NewAccount(newAccountAddress)
			newAccount.SetNonce(newAccount.GetNonce() + 1)
			maybe.PushError(evm.cache.UpdateAccount(newAccount))
			// Run the input to get the contract code.
			// NOTE: no need to copy 'input' as per Call contract.
			// record old ctx
			prevInput := ctx.Input
			prevValue := ctx.Value
			ctx.Input = nil
			ctx.Value = contractValue
			ret, callErr := evm.Call(callee, newAccountAddress, input)
			ctx.Input = prevInput
			ctx.Value = prevValue
			if callErr != nil {
				stack.Push(core.Zero256)
				// Note we both set the return buffer and return the result normally in order to service the error to
				// EVM caller
				returnData = ret
			} else {
				// Update the account with its initialised contract code
				// todo: we may need to set ancestor?
				createDataGas := uint64(len(ret)) * gas.CreateData
				maybe.PushError(useGasNegative(ctx.Gas, createDataGas))
				if maybe.Error() == nil {
					newAccount := evm.getAccount(newAccountAddress)
					newAccount.SetCode(ret)
					stack.PushAddress(newAccountAddress)
				}
			}
			*ctx.Gas += gasPrev

		case CALL, CALLCODE:
			returnData = nil

			var err error
			maybe.PushError(useGasNegative(ctx.Gas, gas.Call))

			var gasLimit = stack.PopUint64()
			gasLimit = callGas(*ctx.Gas, gasLimit)
			maybe.PushError(useGasNegative(ctx.Gas, gasLimit))

			target, value := stack.PopAddress(), stack.PopUint64()
			inOffset, inSize := stack.PopBigInt(), stack.PopBigInt()
			retOffset, retSize := stack.PopBigInt(), stack.PopUint64()
			if value != 0 {
				maybe.PushError(useGasNegative(ctx.Gas, gas.CallValue))
				if op == CALL && isEmptyAccount(evm.getAccount(target)) {
					useGasNegative(ctx.Gas, gas.CallNewAccount)
				}
				gasLimit += gas.CallStipend
			}
			input, memoryGas := memory.Read(inOffset, inSize)
			maybe.PushError(useGasNegative(ctx.Gas, memoryGas))

			// store prev ctx
			prevInput := evm.ctx.Input
			prevValue := evm.ctx.Value
			prevGas := evm.ctx.Gas
			// update ctx
			ctx.Input = input
			ctx.Value = value
			ctx.Gas = &gasLimit
			if debug {
				log.Debugf("  %v", target.Bytes())
			}
			if op == CALL {
				returnData, err = evm.Call(callee, target, evm.getAccount(target).GetCode())
			} else {
				returnData, err = evm.Call(callee, callee, evm.getAccount(target).GetCode())
			}
			if err != nil {
				stack.Push(core.Zero256)
			} else {
				stack.Push(core.One256)
			}
			if err == nil || err.Error() == errors.ExecutionReverted.Error() {
				memory.Write(retOffset, util.RightPadBytes(returnData, int(retSize)))
			}
			// restore ctx
			ctx.Input = prevInput
			ctx.Value = prevValue
			*prevGas += *ctx.Gas
			ctx.Gas = prevGas

		case STATICCALL, DELEGATECALL:
			// todo: support read only mode of STATICCALL
			returnData = nil

			var err error
			maybe.PushError(useGasNegative(ctx.Gas, gas.Call))

			var gas = stack.PopUint64()

			target := stack.PopAddress()
			inOffset, inSize := stack.PopBigInt(), stack.PopBigInt()
			retOffset, retSize := stack.PopBigInt(), stack.PopUint64()
			var memoryGas uint64
			var input []byte
			if op == STATICCALL {
				x, _ := memory.CalMemGas(inOffset.Uint64(), inSize.Uint64())
				y, _ := memory.CalMemGas(retOffset.Uint64(), retSize)
				if x > y {
					memoryGas = x
				} else {
					memoryGas = y
				}
				input, _ = memory.Read(inOffset, inSize)
			} else {
				input, memoryGas = memory.Read(inOffset, inSize)
			}
			gas = staticCallGas(*ctx.Gas, memoryGas, gas)
			maybe.PushError(useGasNegative(ctx.Gas, gas+memoryGas))
			// store prev ctx
			prevInput := evm.ctx.Input
			prevValue := evm.ctx.Value
			prevGas := evm.ctx.Gas
			// update ctx
			ctx.Input = input
			ctx.Value = 0
			ctx.Gas = &gas
			if debug {
				log.Debugf("  %v", target.Bytes())
			}
			if op == STATICCALL {
				returnData, err = evm.CallWithoutTransfer(callee, target, evm.getAccount(target).GetCode())
			} else {
				returnData, err = evm.CallWithoutTransfer(caller, callee, evm.getAccount(target).GetCode())
			}

			if err != nil {
				stack.Push(core.Zero256)
			} else {
				stack.Push(core.One256)
			}
			if err == nil || err.Error() == errors.ExecutionReverted.Error() {
				memory.Write(retOffset, util.RightPadBytes(returnData, int(retSize)))
			}
			// restore ctx
			ctx.Input = prevInput
			ctx.Value = prevValue
			*prevGas += *ctx.Gas
			ctx.Gas = prevGas

		case RETURN: // 0xF3
			maybe.PushError(useGasNegative(ctx.Gas, gas.Zero))
			offset, size := stack.PopBigInt(), stack.PopBigInt()
			output, memoryGas := memory.Read(offset, size)
			maybe.PushError(useGasNegative(ctx.Gas, memoryGas))
			if debug {
				log.Debugf("  [%v, %v] (%d) 0x%X", offset, size, len(output), output)
			}
			return output, maybe.Error()

		case REVERT: // 0xFD
			maybe.PushError(useGasNegative(ctx.Gas, gas.Zero))
			offset, size := stack.PopBigInt(), stack.PopBigInt()
			output, memoryGas := memory.Read(offset, size)
			maybe.PushError(useGasNegative(ctx.Gas, memoryGas))
			if debug {
				log.Debugf("  [%v, %v] (%d) 0x%X", offset, size, len(output), output)
			}
			maybe.PushError(errors.ExecutionReverted)
			return output, maybe.Error()

		case INVALID: // 0xFE
			maybe.PushError(errors.ExecutionAborted)
			useGasNegative(ctx.Gas, *ctx.Gas)
			return nil, maybe.Error()

		case SELFDESTRUCT: // 0xFF
			maybe.PushError(useGasNegative(ctx.Gas, gas.SelfdestructEIP150))
			receiver := stack.PopAddress()
			//todo: different db implementation
			if !evm.cache.Exist(receiver) {
				maybe.PushError(useGasNegative(ctx.Gas, gas.CreateBySelfdestruct))
			}
			account := evm.getAccount(receiver)
			balance := evm.getAccount(callee).GetBalance()
			if isEmptyAccount(account) && balance != 0 {
				maybe.PushError(useGasNegative(ctx.Gas, gas.CreateBySelfdestruct))
			}
			if !evm.cache.HasSuicide(callee) {
				evm.addRefund(gas.SelfdestructRefund)
			}
			maybe.PushError(account.AddBalance(balance))
			maybe.PushError(evm.cache.UpdateAccount(account))
			maybe.PushError(evm.cache.Suicide(callee))
			if debug {
				log.Debugf("  (%v) %v", receiver, balance)
			}
			return nil, maybe.Error()

		case STOP: // 0x00
			maybe.PushError(useGasNegative(ctx.Gas, gas.Zero))
			log.Debugf("")
			return nil, maybe.Error()

		default:
			maybe.PushError(errors.UnknownOpcode)
			if debug {
				log.Debugf("(pc) %-3v Unknown opcode %v", pc, op)
			}
			return nil, maybe.Error()
		}
		pc++
	}
}

// todo: Notice that creator is not used now
func (evm *EVM) createAccount(creator, address Address) error {
	if evm.cache.Exist(address) {
		return errors.InvalidAddress
	}

	account := evm.bc.NewAccount(address)

	return evm.cache.UpdateAccount(account)
}

// todo: if there is a better way to do this?
func getOpCode(code []byte, n uint64) OpCode {
	if uint64(len(code)) <= n {
		return STOP
	}
	return OpCode(code[n])
}

func useGasNegative(gasLeft *uint64, gasToUse uint64) error {
	if *gasLeft >= gasToUse {
		*gasLeft -= gasToUse
	} else {
		return errors.InsufficientGas
	}
	return nil
}

// getAccount is a wrapper of evm.cache.GetAccount
func (evm *EVM) getAccount(address Address) Account {
	return evm.cache.GetAccount(address)
}

func jump(code []byte, to uint64, pc *uint64) error {
	dest := getOpCode(code, to)
	if dest != JUMPDEST {
		log.Debugf("~> %v invalid jump dest %v", to, dest)
		return errors.InvalidJumpDest
	}
	log.Debugf("~> %v", to)
	*pc = to
	return nil
}

func isEmptyValue(bytes []byte) bool {
	if len(bytes) == 0 {
		return true
	}
	for i := range bytes {
		if bytes[i] != 0 {
			return false
		}
	}
	return true
}

// isEqual will compare two bytes
func isEqual(a, b []byte) bool {
	if isEmptyValue(a) && isEmptyValue(b) {
		return true
	}
	return bytes.Equal(a, b)
}

func callGas(availableGas, callCostGas uint64) uint64 {
	availableGas -= 2
	gas := availableGas - availableGas/64
	if gas < callCostGas {
		return gas
	}
	return callCostGas
}

func staticCallGas(availableGas, base, callCost uint64) uint64 {
	availableGas -= base
	availableGas -= availableGas / 64
	if availableGas < callCost {
		return availableGas
	}
	return callCost
}
func isEmptyAccount(account Account) bool {
	if account == nil {
		return true
	}
	if account.GetBalance() == 0 && len(account.GetCode()) == 0 && account.GetNonce() == 0 {
		return true
	}
	return false
}

func wordGas(length, copyGas uint64) uint64 {
	return (length + 31) / 32 * copyGas
}
