package core

import (
	"evm/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWord256Copy(t *testing.T) {
	var word256 Word256
	word256[0] = byte(1)
	copy := word256.Copy()
	require.Equal(t, word256.Bytes(), copy.Bytes())
	word256[0] = byte(2)
	require.NotEqual(t, word256.Bytes(), copy.Bytes())
}

func TestWord256ToWord160(t *testing.T) {
	var word160 Word160
	for i := range word160 {
		word160[i] = byte(util.RandNum(10))
	}
	word256 := word160.Word256()
	require.Equal(t, word160.Bytes(), word256.Word160().Bytes())
}
