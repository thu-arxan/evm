package precompile

import (
	"evm/gas"
	"evm/util"

	"golang.org/x/crypto/ripemd160"
)

// RIPEMD160 implemented as a native contract.
type ripemd160hash struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
//
// This method does not require any overflow checking as the input size gas costs
// required for anything significant is so high it's impossible to pay for.
func (c *ripemd160hash) RequiredGas(input []byte) uint64 {
	return uint64(len(input)+31)/32*gas.Ripemd160PerWord + gas.Ripemd160Base
}
func (c *ripemd160hash) Run(input []byte) ([]byte, error) {
	ripemd := ripemd160.New()
	ripemd.Write(input)
	return util.LeftPadBytes(ripemd.Sum(nil), 32), nil
}