package evm

import (
	"evm/core"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

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
