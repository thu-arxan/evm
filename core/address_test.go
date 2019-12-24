package core

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddressFromBytes(t *testing.T) {
	var bytes = randBytes(20)
	address := AddressFromBytes(bytes)
	require.Equal(t, bytes, address.Bytes())
	bytes = randBytes(12)
	address = AddressFromBytes(bytes)
	require.True(t, strings.HasSuffix(fmt.Sprintf("%x", address.Bytes()), fmt.Sprintf("%x", bytes)))
	bytes = randBytes(32)
	address = AddressFromBytes(bytes)
	require.True(t, strings.HasSuffix(fmt.Sprintf("%x", bytes), fmt.Sprintf("%x", address.Bytes())))
}
