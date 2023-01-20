package datasource

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"

	_ "embed"

	"github.com/NethermindEth/juno/clients"
	"github.com/NethermindEth/juno/core/felt"
)

var (
	//go:embed testdata/class_genesis.json
	classGenesisBytes []byte
	//go:embed testdata/class_0_8.json
	classCairo08Bytes []byte
	//go:embed testdata/class_0_10.json
	classCairo10Bytes []byte
)

func TestProgramHash(t *testing.T) {
	hexToFelt := func(hex string) *felt.Felt {
		f, _ := new(felt.Felt).SetString(hex)
		return f
	}
	tests := []struct {
		class []byte
		want  *felt.Felt
	}{
		{
			class: classGenesisBytes,
			want:  hexToFelt("0x1e87d79be8c8146494b5c54318f7d194481c3959752659a1e1bce158649a670"),
		},
		{
			class: classCairo08Bytes,
			want:  hexToFelt("0x359145fc6207854bfbbeadae4c6e289024400a5af87090ed18073200fef6213"),
		},
		{
			class: classCairo10Bytes,
			want:  hexToFelt("0x88562ac88adfc7760ff452d048d39d72978bcc0f8d7b0fcfb34f33970b3df0"),
		},
	}

	for _, tt := range tests {
		t.Run("ProgramHash", func(t *testing.T) {
			var classDefinition *clients.ClassDefinition
			if err := json.Unmarshal(tt.class, &classDefinition); err != nil {
				t.Fatalf("unexpected error while unmarshaling contract definition: %s", err)
			}

			programHash, err := ProgramHash(classDefinition)
			if err != nil {
				t.Fatalf("unexpected error while computing program hash: %s", err)
			}

			if !programHash.Equal(tt.want) {
				t.Errorf("wrong hash: got %s, want %s", programHash.Text(16), tt.want.Text(16))
			}
		})
	}
}

func TestMarshalJsonWithSeparators(t *testing.T) {
	list := []int{1, 4, 5, 6}
	json, err := MarshalJSONWithSeparators(list, ", ", ": ")
	assert.NoError(t, err)
	assert.Equal(t, `[1, 4, 5, 6]`, string(json))
	/*
		s := struct {
			A int    `json:"asd"`
			B string `json:"123"`
		}{A: 5, B: "444"}
		json, err = MarshalJSONWithSeparators(s, ", ", ": ")
		assert.NoError(t, err)
		assert.Equal(t, `{"123": "444", "asd": 5}`, string(json))
	*/
	hint := clients.Hints{
		2:   "3333",
		100: "asdads",
	}

	json, err = MarshalJSONWithSeparators(struct {
		Hints clients.Hints
	}{Hints: hint}, ", ", ": ")
	assert.NoError(t, err)
	assert.Equal(t, `{"123": "444", "asd": 5}`, string(json))
}
