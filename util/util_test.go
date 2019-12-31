package util

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLog256(t *testing.T) {
	require.Equal(t, 0, Log256(big.NewInt(1)))
	require.Equal(t, 0, Log256(big.NewInt(255)))
	require.Equal(t, 1, Log256(big.NewInt(256)))
	require.Equal(t, 1, Log256(big.NewInt(65535)))
	require.Equal(t, 2, Log256(big.NewInt(65536)))
}

func TestFixBytesLength(t *testing.T) {
	require.Len(t, FixBytesLength(nil, 32), 32)
}
