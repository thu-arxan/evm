package evm

import (
	"bytes"
	"math/big"

	"evm/core"
	"evm/errors"
	"evm/gas"
	"evm/precompile"
	"evm/util"

	"evm/crypto"

	"github.com/labstack/gommon/log"
)

// Here defines some default stack capacity variables
const (
	DefaultStackCapacity    uint64 = 1024
	DefaultMaxStackCapacity uint64 = 32 * 1024
	MaxCodeSize             int    = 24576
)

func init() {
	log.SetLevel(log.DEBUG)
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
}

// New is the constructor of EVM
func New(bc Blockchain, db DB, ctx *Context) *EVM {
	return &EVM{
		bc:             bc,
		cache:          NewCache(db),
		memoryProvider: DefaultDynamicMemoryProvider,
		ctx:            ctx,
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

	evm.cache.Sync()
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
	evm.cache.Sync()
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
			return nil, maybe.Error()
		}

		var op = getOpCode(code, pc)
		log.Debugf("(pc) %-3d (op) %-14s (st) %-4d (gas) %d", pc, op.String(), stack.Len(), *ctx.Gas)

		switch op {
		case ADD: // 0x01
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopBigInt(), stack.PopBigInt()
			sum := new(big.Int).Add(x, y)
			res := stack.PushBigInt(sum)
			log.Debugf("%v + %v = %v (%v)", x, y, sum, res)

		case MUL: // 0x02
			maybe.PushError(useGasNegative(ctx.Gas, gas.Low))
			x, y := stack.PopBigInt(), stack.PopBigInt()
			prod := new(big.Int).Mul(x, y)
			res := stack.PushBigInt(prod)
			log.Debugf("%v * %v = %v (%v)", x, y, prod, res)

		case SUB: // 0x03
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopBigInt(), stack.PopBigInt()
			diff := new(big.Int).Sub(x, y)
			res := stack.PushBigInt(diff)
			log.Debugf("%v - %v = %v (%v)", x, y, diff, res)

		case DIV: // 0x04
			maybe.PushError(useGasNegative(ctx.Gas, gas.Low))
			x, y := stack.PopBigInt(), stack.PopBigInt()
			if y.Sign() == 0 {
				stack.Push(core.Zero256)
				log.Debugf("%v / %v = %v", x, y, 0)
			} else {
				div := new(big.Int).Div(x, y)
				res := stack.PushBigInt(div)
				log.Debugf("%v / %v = %v (%v)", x, y, div, res)
			}

		case SDIV: // 0x05
			maybe.PushError(useGasNegative(ctx.Gas, gas.Low))
			x, y := stack.PopSignedBigInt(), stack.PopSignedBigInt()
			if x.Sign() == 0 || y.Sign() == 0 {
				stack.Push(core.Zero256)
				log.Debugf("%v / %v = %v", x, y, 0)
			} else {
				var div *big.Int
				if x.Sign() != y.Sign() {
					div = new(big.Int).Div(x.Abs(x), y.Abs(y))
					div.Neg(div)
				} else {
					div = new(big.Int).Div(x.Abs(x), y.Abs(y))
				}
				res := stack.PushBigInt(div)
				log.Debugf("%v / %v = %v (%v)", x, y, div, res)
			}

		case MOD: // 0x06
			maybe.PushError(useGasNegative(ctx.Gas, gas.Low))
			x, y := stack.PopBigInt(), stack.PopBigInt()
			if y.Sign() == 0 {
				stack.Push(core.Zero256)
				log.Debugf("%v %% %v = %v", x, y, 0)
			} else {
				mod := new(big.Int).Mod(x, y)
				res := stack.PushBigInt(mod)
				log.Debugf("%v %% %v = %v (%v)", x, y, mod, res)
			}

		case SMOD: // 0x07
			maybe.PushError(useGasNegative(ctx.Gas, gas.Low))
			x, y := stack.PopSignedBigInt(), stack.PopSignedBigInt()
			if y.Sign() == 0 {
				stack.Push(core.Zero256)
				log.Debugf("%v %% %v = %v", x, y, 0)
			} else {
				mod := new(big.Int).Mod(x, y)
				res := stack.PushBigInt(mod)
				log.Debugf("%v %% %v = %v (%v)", x, y, mod, res)
			}

		case ADDMOD: // 0x08
			maybe.PushError(useGasNegative(ctx.Gas, gas.Mid))
			x, y, z := stack.PopBigInt(), stack.PopBigInt(), stack.PopBigInt()
			if z.Sign() == 0 {
				stack.Push(core.Zero256)
				log.Debugf("%v %% %v = %v\n", x, y, 0)
			} else {
				add := new(big.Int).Add(x, y)
				mod := add.Mod(add, z)
				res := stack.PushBigInt(mod)
				log.Debugf("%v + %v %% %v = %v (%v)\n", x, y, z, mod, res)
			}

		case MULMOD: // 0x09
			maybe.PushError(useGasNegative(ctx.Gas, gas.Mid))
			x, y, z := stack.PopBigInt(), stack.PopBigInt(), stack.PopBigInt()
			if z.Sign() == 0 {
				stack.Push(core.Zero256)
				log.Debugf("%v %% %v = %v\n", x, y, 0)
			} else {
				mul := new(big.Int).Mul(x, y)
				mod := mul.Mod(mul, z)
				res := stack.PushBigInt(mod)
				log.Debugf("%v * %v %% %v = %v (%v)\n", x, y, z, mod, res)
			}

		case EXP: // 0x0A
			x, y := stack.PopBigInt(), stack.PopBigInt()
			if y.Sign() == 0 {
				maybe.PushError(useGasNegative(ctx.Gas, gas.Exp))
			} else {
				maybe.PushError(useGasNegative(ctx.Gas, gas.Exp+gas.ExpByte*uint64(1+util.Log256(y))))
			}
			pow := new(big.Int).Exp(x, y, nil)
			res := stack.PushBigInt(pow)
			log.Debugf("%v ** %v = %v (%v)\n", x, y, pow, res)

		case SIGNEXTEND: // 0x0B
			maybe.PushError(useGasNegative(ctx.Gas, gas.Low))
			back := stack.PopUint64()
			if back < core.Word256Bytes-1 {
				bits := uint(back*8 + 7)
				stack.PushBigInt(core.SignExtend(stack.PopBigInt(), bits))
			}

		case LT: // 0x10
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopBigInt(), stack.PopBigInt()
			if x.Cmp(y) < 0 {
				stack.Push(core.One256)
				log.Debugf("%v < %v = %v\n", x, y, 1)
			} else {
				stack.Push(core.Zero256)
				log.Debugf("%v < %v = %v\n", x, y, 0)
			}

		case GT: // 0x11
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopBigInt(), stack.PopBigInt()
			if x.Cmp(y) > 0 {
				stack.Push(core.One256)
				log.Debugf("%v > %v = %v\n", x, y, 1)
			} else {
				stack.Push(core.Zero256)
				log.Debugf("%v > %v = %v\n", x, y, 0)
			}

		case SLT: // 0x12
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopSignedBigInt(), stack.PopSignedBigInt()
			if x.Cmp(y) < 0 {
				stack.Push(core.One256)
				log.Debugf("%v < %v = %v\n", x, y, 1)
			} else {
				stack.Push(core.Zero256)
				log.Debugf("%v < %v = %v\n", x, y, 0)
			}

		case SGT: // 0x13
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.PopSignedBigInt(), stack.PopSignedBigInt()
			if x.Cmp(y) > 0 {
				stack.Push(core.One256)
				log.Debugf("%v > %v = %v\n", x, y, 1)
			} else {
				stack.Push(core.Zero256)
				log.Debugf("%v > %v = %v\n", x, y, 0)
			}

		case EQ: // 0x14
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.Pop(), stack.Pop()
			if isEqual(x[:], y[:]) {
				stack.Push(core.One256)
				log.Debugf("%v == %v = %v\n", x, y, 1)
			} else {
				stack.Push(core.Zero256)
				log.Debugf("%v == %v = %v\n", x, y, 0)
			}

		case ISZERO: // 0x15
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x := stack.Pop()
			if x.IsZero() {
				stack.Push(core.One256)
				log.Debugf("%v == 0 = %v\n", x, 1)
			} else {
				stack.Push(core.Zero256)
				log.Debugf("%v == 0 = %v\n", x, 0)
			}

		case AND: // 0x16
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.Pop(), stack.Pop()
			z := [32]byte{}
			for i := 0; i < 32; i++ {
				z[i] = x[i] & y[i]
			}
			stack.Push(z)
			log.Debugf("%v & %v = %v\n", x, y, z)

		case OR: // 0x17
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.Pop(), stack.Pop()
			z := [32]byte{}
			for i := 0; i < 32; i++ {
				z[i] = x[i] | y[i]
			}
			stack.Push(z)
			log.Debugf("%v | %v = %v\n", x, y, z)

		case XOR: // 0x18
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x, y := stack.Pop(), stack.Pop()
			z := [32]byte{}
			for i := 0; i < 32; i++ {
				z[i] = x[i] ^ y[i]
			}
			stack.Push(z)
			log.Debugf("%v ^ %v = %v\n", x, y, z)

		case NOT: // 0x19
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			x := stack.Pop()
			z := [32]byte{}
			for i := 0; i < 32; i++ {
				z[i] = ^x[i]
			}
			stack.Push(z)
			log.Debugf("!%v = %v\n", x, z)

		case BYTE: // 0x1A
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			idx := stack.PopUint64()
			val := stack.Pop()
			res := byte(0)
			if idx < 32 {
				res = val[idx]
			}
			stack.PushUint64(uint64(res))
			log.Debugf("=> 0x%X\n", res)

		case SHL: //0x1B
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			shift, x := stack.PopBigInt(), stack.PopBigInt()
			if shift.Cmp(core.Big256) >= 0 {
				reset := big.NewInt(0)
				stack.PushBigInt(reset)
				log.Debugf("%v << %v = %v\n", x, shift, reset)
			} else {
				shiftedValue := x.Lsh(x, uint(shift.Uint64()))
				stack.PushBigInt(shiftedValue)
				log.Debugf("%v << %v = %v\n", x, shift, shiftedValue)
			}

		case SHR: //0x1C
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			shift, x := stack.PopBigInt(), stack.PopBigInt()
			if shift.Cmp(core.Big256) >= 0 {
				reset := big.NewInt(0)
				stack.PushBigInt(reset)
				log.Debugf("%v << %v = %v\n", x, shift, reset)
			} else {
				shiftedValue := x.Rsh(x, uint(shift.Uint64()))
				stack.PushBigInt(shiftedValue)
				log.Debugf("%v << %v = %v\n", x, shift, shiftedValue)
			}

		case SAR: //0x1D
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			shift, x := stack.PopBigInt(), stack.PopSignedBigInt()
			if shift.Cmp(core.Big256) >= 0 {
				reset := big.NewInt(0)
				if x.Sign() < 0 {
					reset.SetInt64(-1)
				}
				stack.PushBigInt(reset)
				log.Debugf("%v << %v = %v\n", x, shift, reset)
			} else {
				shiftedValue := x.Rsh(x, uint(shift.Uint64()))
				stack.PushBigInt(shiftedValue)
				log.Debugf("%v << %v = %v\n", x, shift, shiftedValue)
			}

		case SHA3: // 0x20
			offset, size := stack.PopBigInt(), stack.PopBigInt()
			maybe.PushError(useGasNegative(ctx.Gas, gas.SHA3+gas.SHA3Word*((size.Uint64()+31)/32)))
			data, memoryGas := memory.Read(offset, size)
			maybe.PushError(useGasNegative(ctx.Gas, memoryGas))
			data = crypto.Keccak256(data)
			stack.PushBytes(data)
			log.Debugf("=> (%v) %X\n", size, data)

		case ADDRESS: // 0x30
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushAddress(callee)
			log.Debugf("=> %v\n", callee)

		case BALANCE: // 0x31
			maybe.PushError(useGasNegative(ctx.Gas, gas.Balance))
			address := stack.PopAddress()
			balance := evm.getAccount(address).GetBalance()
			stack.PushUint64(balance)
			log.Debugf("=> %v (%v)\n", balance, address)

		case ORIGIN: // 0x32
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushAddress(evm.origin)
			log.Debugf("=> %v\n", evm.origin)

		case CALLER: // 0x33
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushAddress(caller)
			log.Debugf("=> %v\n", caller)

		case CALLVALUE: // 0x34
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushUint64(ctx.Value)
			log.Debugf("=> %v\n", ctx.Value)

		case CALLDATALOAD: // 0x35
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			offset := stack.PopUint64()
			data, err := util.SubSlice(ctx.Input, offset, 32)
			if err != nil {
				maybe.PushError(errors.InputOutOfBounds)
			}
			res := core.LeftPadWord256(data)
			stack.Push(res)
			log.Debugf("=> 0x%v\n", res)

		case CALLDATASIZE: // 0x36
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushUint64(uint64(len(ctx.Input)))
			log.Debugf("=> %d\n", len(ctx.Input))

		case CALLDATACOPY: // 0x37
			memOff := stack.PopBigInt()
			inputOff := stack.PopUint64()
			length := stack.PopUint64()
			data, err := util.SubSlice(ctx.Input, inputOff, length)
			if err != nil {
				maybe.PushError(errors.InputOutOfBounds)
			}
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow+gas.Copy*((length+31)/32)))
			gasCost := memory.Write(memOff, data)
			maybe.PushError(useGasNegative(ctx.Gas, gasCost))
			log.Debugf("=> [%v, %v, %v] %X\n", memOff, inputOff, length, data)

		case CODESIZE: // 0x38
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			l := uint64(len(code))
			stack.PushUint64(l)
			log.Debugf("=> %d\n", l)

		case CODECOPY: // 0x39
			memOff := stack.PopBigInt()
			codeOff := stack.PopUint64()
			length := stack.PopUint64()
			data, err := util.SubSlice(code, codeOff, length)
			if err != nil {
				maybe.PushError(errors.InputOutOfBounds)
			}
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow+gas.Copy*((length+31)/32)))
			gasCost := memory.Write(memOff, data)
			maybe.PushError(useGasNegative(ctx.Gas, gasCost))
			log.Debugf("=> [%v, %v, %v] %X\n", memOff, codeOff, length, data)

		case GASPRICE: // 0x3A
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushUint64(ctx.GasPrice)
			log.Debugf("=> %v\n", ctx.GasPrice)

		case EXTCODESIZE: // 0x3B
			maybe.PushError(useGasNegative(ctx.Gas, gas.ExtCode))
			address := stack.PopAddress()
			acc := evm.getAccount(address)
			length := uint64(len(acc.GetCode()))
			stack.PushUint64(length)
			log.Debugf("=> %d\n", length)

		case EXTCODECOPY: // 0x3C
			address := stack.PopAddress()
			code := evm.getAccount(address).GetCode()
			memOff := stack.PopBigInt()
			codeOff := stack.PopUint64()
			length := stack.PopUint64()
			data, err := util.SubSlice(code, codeOff, length)
			if err != nil {
				maybe.PushError(errors.InputOutOfBounds)
			}
			gasCost := memory.Write(memOff, data)
			maybe.PushError(useGasNegative(ctx.Gas, gas.ExtCode+gasCost))
			log.Debugf("=> [%v, %v, %v] %X\n", memOff, codeOff, length, data)

		case RETURNDATASIZE: // 0x3D
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushUint64(uint64(len(returnData)))
			log.Debugf("=> %d\n", len(returnData))

		case RETURNDATACOPY: // 0x3E
			memOff, outputOff, length := stack.PopBigInt(), stack.PopBigInt(), stack.PopBigInt()
			end := new(big.Int).Add(outputOff, length)

			if end.BitLen() > 64 || uint64(len(returnData)) < end.Uint64() {
				maybe.PushError(errors.ReturnDataOutOfBounds)
				continue
			}
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow+gas.Copy*((length.Uint64()+31)/32)))
			gasCost := memory.Write(memOff, returnData)
			maybe.PushError(useGasNegative(ctx.Gas, gasCost))
			log.Debugf("=> [%v, %v, %v] %X\n", memOff, outputOff, length, returnData)

		case EXTCODEHASH: // 0x3F
			maybe.PushError(useGasNegative(ctx.Gas, gas.ExtcodeHash))
			address := stack.PopAddress()
			acc := evm.getAccount(address)
			// keccak256 hash of a contract's code
			var extcodehash core.Word256
			if len(acc.GetCodeHash()) > 0 {
				copy(extcodehash[:], acc.GetCodeHash())
			} else {
				copy(extcodehash[:], crypto.Keccak256(acc.GetCode()))
			}
			stack.Push(extcodehash)

		case BLOCKHASH: // 0x40
			maybe.PushError(useGasNegative(ctx.Gas, gas.BlockHash))
			blockNumber := stack.PopUint64()
			// Note: Here is >= other than > because block is not generated while running tx
			if blockNumber >= ctx.BlockHeight {
				log.Debugf("=> attempted to get block hash of a non-existent block: %v", blockNumber)
				maybe.PushError(errors.InvalidBlockNumber)
			} else if ctx.BlockHeight-blockNumber > 256 {
				log.Debugf("=> attempted to get block hash of a block %d outof range", blockNumber)
				maybe.PushError(errors.BlockNumberOutOfRange)
			} else {
				blockHash, err := evm.bc.GetBlockHash(blockNumber)
				if err != nil {
					maybe.PushError(err)
				}
				stack.Push(core.LeftPadWord256(blockHash))
				log.Debugf("=> 0x%v\n", blockHash)
			}

		case COINBASE: // 0x41
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushBytes(ctx.CoinBase)
			log.Debugf("=> 0x%v (NOT SUPPORTED)\n", stack.Peek())

		case TIMESTAMP: // 0x42
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			blockTime := ctx.BlockTime
			stack.PushUint64(uint64(blockTime))
			log.Debugf("=> %d\n", blockTime)

		case NUMBER: // 0x43
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			number := ctx.BlockHeight
			stack.PushUint64(number)
			log.Debugf("=> %d\n", number)

		case DIFFICULTY: // Note: New version deprecated
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			difficulty := ctx.Difficulty
			stack.PushUint64(difficulty)
			log.Debugf("=> %d\n", difficulty)

		case GASLIMIT: // 0x45
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushUint64(ctx.GasLimit)
			log.Debugf("=> %v\n", ctx.GasLimit)

		case CHAINID: // 0x46
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.Push(core.Word256{})
			log.Debugf("Not implemented")

		case SELFBALANCE: // 0x47
			maybe.PushError(useGasNegative(ctx.Gas, gas.Low))
			balance := evm.getAccount(callee).GetBalance()
			stack.PushUint64(balance)
			log.Debugf("=> %v (%v)\n", balance, callee)

		case POP: // 0x50
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			popped := stack.Pop()
			log.Debugf("=> 0x%v\n", popped)

		case MLOAD: // 0x51
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			offset := stack.PopBigInt()
			data, memoryGas := memory.Read(offset, core.BigWord256Bytes)
			maybe.PushError(useGasNegative(ctx.Gas, memoryGas))
			stack.Push(core.LeftPadWord256(data))
			log.Debugf("=> 0x%X @ 0x%v\n", data, offset)

		case MSTORE: // 0x52
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			offset, data := stack.PopBigInt(), stack.Pop()
			gasCost := memory.Write(offset, data.Bytes())
			maybe.PushError(useGasNegative(ctx.Gas, gasCost))
			log.Debugf("=> 0x%v @ 0x%v\n", data, offset)

		case MSTORE8: // 0x53
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			offset := stack.PopBigInt()
			val64 := stack.PopUint64()
			val := byte(val64 & 0xFF)
			gasCost := memory.Write(offset, []byte{val})
			maybe.PushError(useGasNegative(ctx.Gas, gasCost))
			log.Debugf("=> [%v] 0x%X\n", offset, val)

		case SLOAD: // 0x54
			maybe.PushError(useGasNegative(ctx.Gas, gas.Sload))
			loc := stack.Pop()
			value := evm.cache.GetStorage(callee, loc)
			data := core.LeftPadWord256(value)
			stack.Push(data)
			log.Debugf("%v {0x%v = 0x%v}\n", callee, loc, data)

		case SSTORE: // 0x55
			loc, data := stack.Pop(), stack.Pop()
			currentData := evm.cache.GetStorage(callee, loc)
			if *ctx.Gas <= gas.SstoreSentryEIP2200 {
				maybe.PushError(errors.InsufficientGas)
			}
			if isEqual(data.Bytes(), currentData) {
				maybe.PushError(useGasNegative(ctx.Gas, gas.SstoreNoopEIP2200))
			} else {
				originData := evm.cache.db.GetStorage(callee, loc)
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
			log.Debugf("%v {%v := %v}\n", callee, loc, data)

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
				log.Debugf("~> false\n")
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
			log.Debugf("=> 0x%X\n", capacity)

		case GAS: // 0x5A
			maybe.PushError(useGasNegative(ctx.Gas, gas.Base))
			stack.PushUint64(*ctx.Gas)
			log.Debugf("=> %X\n", *ctx.Gas)

		case JUMPDEST: // 0x5B
			maybe.PushError(useGasNegative(ctx.Gas, gas.JumpDest))
			log.Debugf("\n")
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
			log.Debugf("=> 0x%v\n", res)

		case DUP1, DUP2, DUP3, DUP4, DUP5, DUP6, DUP7, DUP8, DUP9, DUP10, DUP11, DUP12, DUP13, DUP14, DUP15, DUP16:
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			n := int(op - DUP1 + 1)
			stack.Dup(n)
			log.Debugf("=> [%d] 0x%v\n", n, stack.Peek())

		case SWAP1, SWAP2, SWAP3, SWAP4, SWAP5, SWAP6, SWAP7, SWAP8, SWAP9, SWAP10, SWAP11, SWAP12, SWAP13, SWAP14, SWAP15, SWAP16:
			maybe.PushError(useGasNegative(ctx.Gas, gas.VeryLow))
			n := int(op - SWAP1 + 2)
			stack.Swap(n)
			log.Debugf("=> [%d] %v\n", n, stack.Peek())

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
			log.Debugf("=> T:%v D:%X\n", topics, data)

		case CREATE, CREATE2: // 0xF0, 0xFB
			if err := useGasNegative(ctx.Gas, gas.Create); err != nil {
				return nil, err
			}
			returnData = nil
			contractValue := stack.PopUint64()
			offset, size := stack.PopBigInt(), stack.PopBigInt()
			input, memoryGas := memory.Read(offset, size)
			if op == CREATE2 {
				//TODO: consider overflow
				wordGas := (size.Uint64() + 31) / 32 * gas.SHA3Word
				memoryGas += wordGas
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
				*ctx.Gas += gasPrev
			}

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
			log.Debugf("=> %v\n", target.Bytes())
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
			gas = callGas(*ctx.Gas, gas)

			target := stack.PopAddress()
			inOffset, inSize := stack.PopBigInt(), stack.PopBigInt()
			retOffset, retSize := stack.PopBigInt(), stack.PopUint64()
			input, memoryGas := memory.Read(inOffset, inSize)
			maybe.PushError(useGasNegative(ctx.Gas, memoryGas))
			maybe.PushError(useGasNegative(ctx.Gas, gas))
			// store prev ctx
			prevInput := evm.ctx.Input
			prevValue := evm.ctx.Value
			prevGas := evm.ctx.Gas
			// update ctx
			ctx.Input = input
			ctx.Value = 0
			ctx.Gas = &gas
			log.Debugf("=> %v\n", target.Bytes())
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
			log.Debugf("=> [%v, %v] (%d) 0x%X\n", offset, size, len(output), output)
			return output, maybe.Error()

		case REVERT: // 0xFD
			maybe.PushError(useGasNegative(ctx.Gas, gas.Zero))
			offset, size := stack.PopBigInt(), stack.PopBigInt()
			output, memoryGas := memory.Read(offset, size)
			maybe.PushError(useGasNegative(ctx.Gas, memoryGas))
			log.Debugf("=> [%v, %v] (%d) 0x%X\n", offset, size, len(output), output)
			maybe.PushError(errors.ExecutionReverted)
			return output, maybe.Error()

		case INVALID: // 0xFE
			maybe.PushError(errors.ExecutionAborted)
			return nil, maybe.Error()

		case SELFDESTRUCT: // 0xFF
			maybe.PushError(useGasNegative(ctx.Gas, gas.SelfdestructEIP150))
			receiver := stack.PopAddress()
			account := evm.getAccount(receiver)
			balance := evm.getAccount(callee).GetBalance()
			if isEmptyAccount(account) && balance != 0 {
				maybe.PushError(useGasNegative(ctx.Gas, gas.CreateBySelfdestruct))
			}
			if evm.cache.HasSuicide(callee) {
				evm.addRefund(gas.SelfdestructRefund)
			}
			maybe.PushError(account.AddBalance(balance))
			maybe.PushError(evm.cache.UpdateAccount(account))
			maybe.PushError(evm.cache.Suicide(callee))
			log.Debugf("=> (%v) %v\n", receiver, balance)
			return nil, maybe.Error()

		case STOP: // 0x00
			maybe.PushError(useGasNegative(ctx.Gas, gas.Zero))
			log.Debugf("\n")
			return nil, maybe.Error()

		default:
			maybe.PushError(errors.UnknownOpcode)
			log.Debugf("(pc) %-3v Unknown opcode %v\n", pc, op)
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
		log.Debugf("~> %v invalid jump dest %v\n", to, dest)
		return errors.InvalidJumpDest
	}
	log.Debugf("~> %v\n", to)
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

func isEmptyAccount(account Account) bool {
	if account == nil {
		return true
	}
	if account.GetBalance() == 0 && len(account.GetCode()) == 0 && account.GetNonce() == 0 {
		return true
	}
	return false
}
