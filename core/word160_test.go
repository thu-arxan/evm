package core

import (
	"evm/util"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWord160ToWord256(t *testing.T) {
	var word160 Word160
	for i := range word160 {
		word160[i] = byte(util.RandNum(10))
	}
	word160Hex := fmt.Sprintf("%x", word160.Bytes())
	word256Hex := fmt.Sprintf("%x", word160.Word256().Bytes())
	require.Equal(t, "000000000000000000000000"+word160Hex, word256Hex)
}
