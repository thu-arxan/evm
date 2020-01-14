package evm

import (
	"evm/core"
	"evm/util"
	"fmt"
	"math"
	"math/big"
	"time"

	"evm/errors"
)

// Here defines some cost
var (
	PushTime       int64
	PushSize       int64
	PushUintTime   int64
	PushUintSize   int64
	PushBigIntTime int64
	PushBigIntSize int64
	PopTime        int64
	PopSize        int64
	PopUintTime    int64
	PopUintSize    int64
	PopBigIntTime  int64
	PopBigIntSize  int64
)

// Stack is the stack that support the running of evm
// Note: The stack is not thread safety
type Stack struct {
	data        []core.Word256
	maxCapacity uint64
	ptr         int

	gas     *uint64
	errSink errors.Sink

	toAddressFunc func(bytes []byte) Address
}

// NewStack is the constructor of Stack
func NewStack(initialCapacity uint64, maxCapacity uint64, gas *uint64, errSink errors.Sink, toAddressFunc func(bytes []byte) Address) *Stack {
	return &Stack{
		data:          make([]core.Word256, initialCapacity),
		maxCapacity:   maxCapacity,
		gas:           gas,
		errSink:       errSink,
		toAddressFunc: toAddressFunc,
	}
}

// Push push core.Word256 into stack
func (st *Stack) Push(word core.Word256) {
	now := time.Now()
	err := st.ensureCapacity(uint64(st.ptr) + 1)
	if err != nil {
		st.pushErr(errors.DataStackOverflow)
		return
	}
	st.data[st.ptr] = word
	st.ptr++
	PushSize++
	PushTime += int64(time.Since(now))
}

// Pop pos a core.Word256 from the stak
func (st *Stack) Pop() core.Word256 {
	if st.ptr == 0 {
		st.pushErr(errors.DataStackUnderflow)
		return core.Zero256
	}
	st.ptr--
	return st.data[st.ptr]
}

// PushBytes push bytes into stack, bytes length would fixed to 32
func (st *Stack) PushBytes(bz []byte) {
	bz = util.FixBytesLength(bz, 32)
	st.Push(core.LeftPadWord256(bz))
}

// PushAddress push address into stack
func (st *Stack) PushAddress(address Address) {
	st.Push(core.BytesToWord256(address.Bytes()))
}

// PushUint64 push uint64 into stack
func (st *Stack) PushUint64(i uint64) {
	now := time.Now()
	st.Push(core.Uint64ToWord256(i))
	PushUintSize++
	PushUintTime += int64(time.Since(now))
}

// PopUint64 pop uint64 from stack
func (st *Stack) PopUint64() uint64 {
	word := st.Pop()
	if Is64BitOverflow(word) {
		st.pushErr(fmt.Errorf("uint64 overflow from word: %v", word))
		return 0
	}
	return core.Uint64FromWord256(word)
}

// PushBigInt push the bigInt as a core.Word256 encoding negative values in 32-byte twos complement and returns the encoded result
func (st *Stack) PushBigInt(bigInt *big.Int) core.Word256 {
	now := time.Now()
	word := core.LeftPadWord256(core.U256(bigInt).Bytes())
	st.Push(word)
	PushBigIntSize++
	PushBigIntTime += int64(time.Since(now))
	return word
}

// PopSignedBigInt pop signed big int from stack
func (st *Stack) PopSignedBigInt() *big.Int {
	return core.S256(st.PopBigInt())
}

// PopBigInt pop big int from stack
func (st *Stack) PopBigInt() *big.Int {
	now := time.Now()
	word := st.Pop()
	i := new(big.Int).SetBytes(word[:])
	PopBigIntSize++
	PopBigIntTime += int64(time.Since(now))
	return i
}

// PopBytes pop bytes from stack
func (st *Stack) PopBytes() []byte {
	return st.Pop().Bytes()
}

// PopAddress pop address from stack
func (st *Stack) PopAddress() Address {
	if st.toAddressFunc != nil {
		return st.toAddressFunc(st.Pop().Bytes())
	}
	return st.Pop().Address()
}

// Len return length of stack
func (st *Stack) Len() int {
	return st.ptr
}

// Swap swap stack
func (st *Stack) Swap(n int) {
	// st.useGas(GasStackOp)
	if st.ptr < n {
		st.pushErr(errors.DataStackUnderflow)
		return
	}
	st.data[st.ptr-n], st.data[st.ptr-1] = st.data[st.ptr-1], st.data[st.ptr-n]
}

// Dup duplicate stack
func (st *Stack) Dup(n int) {
	// st.useGas(GasStackOp)
	if st.ptr < n {
		st.pushErr(errors.DataStackUnderflow)
		return
	}
	st.Push(st.data[st.ptr-n])
}

// Peek peek the stack element
func (st *Stack) Peek() core.Word256 {
	if st.ptr == 0 {
		st.pushErr(errors.DataStackUnderflow)
		return core.Zero256
	}
	return st.data[st.ptr-1]
}

// Print print stack status
func (st *Stack) Print(n int) {
	fmt.Println("### stack ###")
	if st.ptr > 0 {
		nn := n
		if st.ptr < n {
			nn = st.ptr
		}
		for j, i := 0, st.ptr-1; i > st.ptr-1-nn; i-- {
			fmt.Printf("%-3d  %X\n", j, st.data[i])
			j++
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("#############")
}

// Is64BitOverflow return if the word overflow
func Is64BitOverflow(word core.Word256) bool {
	for i := 0; i < len(word)-8; i++ {
		if word[i] != 0 {
			return true
		}
	}
	return false
}

// Ensures the current stack can hold a new element. Will only grow the
// backing array (will not shrink).
func (st *Stack) ensureCapacity(newCapacity uint64) error {
	// Maximum length of a data that allocates memory is the same as the native int max size
	// We could rethink this limit, but we don't want different validators to disagree on
	// transaction validity so we pick the lowest common denominator
	if newCapacity > math.MaxInt32 {
		// If we ever did want more than an int32 of space then we would need to
		// maintain multiple pages of memory
		return fmt.Errorf("cannot address memory beyond a maximum index "+"with int32 width (%v bytes)", math.MaxInt32)
	}
	newCapacityInt := int(newCapacity)
	// We're already big enough so return
	if newCapacityInt <= len(st.data) {
		return nil
	}
	if st.maxCapacity > 0 && newCapacity > st.maxCapacity {
		return fmt.Errorf("cannot grow memory because it would exceed the "+
			"current maximum limit of %v bytes", st.maxCapacity)
	}
	// Ensure the backing array of data is big enough
	// Grow the memory one word at time using the pre-allocated zeroWords to avoid
	// unnecessary allocations. Use append to make use of any spare capacity in
	// the data's backing array.
	for newCapacityInt > cap(st.data) {
		// We'll trust Go exponentially grow our arrays (at first).
		st.data = append(st.data, core.Zero256)
	}
	// Now we've ensured the backing array of the data is big enough we can
	// just re-data (even if len(mem.data) < newCapacity)
	st.data = st.data[:newCapacity]
	return nil
}

// func (st *Stack) useGas(gasToUse uint64) {
// 	if *st.gas > gasToUse {
// 		*st.gas -= gasToUse
// 	} else {
// 		st.pushErr(errors.InsufficientGas)
// 	}
// }

func (st *Stack) pushErr(err error) {
	st.errSink.PushError(err)
}
