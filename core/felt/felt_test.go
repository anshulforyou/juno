package felt

import (
	"testing"

	"github.com/NethermindEth/juno/encoder"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJson(t *testing.T) {
	var with Felt
	assert.NoError(t, with.UnmarshalJSON([]byte("0x4437ab")))

	var without Felt
	assert.NoError(t, without.UnmarshalJSON([]byte("4437ab")))
	assert.Equal(t, true, without.Equal(&with))
}

func TestFeltCbor(t *testing.T) {
	var val Felt
	_, err := val.SetRandom()
	assert.NoError(t, err)

	encoder.TestSymmetry(t, val)
}
