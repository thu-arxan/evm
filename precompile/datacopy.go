package precompile

import "evm/gas"

// data copy implemented as a native contract.
type dataCopy struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
//
// This method does not require any overflow checking as the input size gas costs
// required for anything significant is so high it's impossible to pay for.
func (c *dataCopy) RequiredGas(input []byte) uint64 {
	return uint64(len(input)+31)/32*gas.IdentityPerWord + gas.IdentityBase
}
func (c *dataCopy) Run(in []byte) ([]byte, error) {
	return in, nil
}
