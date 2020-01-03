package precompile

import (
	"errors"
	"evm/gas"
)

// ECRECOVER implemented as a native contract.
type ecrecover struct{}

func (c *ecrecover) RequiredGas(input []byte) uint64 {
	return gas.Ecrecover
}

func (c *ecrecover) Run(input []byte) ([]byte, error) {
	return nil, errors.New("Not implementation yet")
	// const ecRecoverInputLength = 128

	// input = common.RightPadBytes(input, ecRecoverInputLength)
	// // "input" is (hash, v, r, s), each 32 bytes
	// // but for ecrecover we want (r, s, v)

	// r := new(big.Int).SetBytes(input[64:96])
	// s := new(big.Int).SetBytes(input[96:128])
	// v := input[63] - 27

	// // tighter sig s values input homestead only apply to tx sigs
	// if !allZero(input[32:63]) || !crypto.ValidateSignatureValues(v, r, s, false) {
	// 	return nil, nil
	// }
	// // We must make sure not to modify the 'input', so placing the 'v' along with
	// // the signature needs to be done on a new allocation
	// sig := make([]byte, 65)
	// copy(sig, input[64:128])
	// sig[64] = v
	// // v needs to be at the end for libsecp256k1
	// pubKey, err := crypto.Ecrecover(input[:32], sig)
	// // make sure the public key is a valid one
	// if err != nil {
	// 	return nil, nil
	// }

	// // the first byte of pubkey is bitcoin heritage
	// return common.LeftPadBytes(crypto.Keccak256(pubKey[1:])[12:], 32), nil
}
