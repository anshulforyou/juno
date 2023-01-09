package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/NethermindEth/juno/core/crypto"
	"github.com/NethermindEth/juno/core/felt"
)

// Class unambiguously defines a [Contract]'s semantics.
type Class struct {
	// The version of the class, currently always 0.
	APIVersion *felt.Felt
	// External functions defined in the class.
	Externals []EntryPoint
	// Functions that receive L1 messages. See
	// https://www.cairo-lang.org/docs/hello_starknet/l1l2.html#receiving-a-message-from-l1
	L1Handlers []EntryPoint
	// Constructors for the class. Currently, only one is allowed.
	Constructors []EntryPoint
	// An ascii-encoded array of builtin names imported by the class.
	Builtins []*felt.Felt
	// The starknet_keccak hash of the ".json" file compiler output.
	ProgramHash *felt.Felt
	Bytecode    []*felt.Felt
}

func (c *Class) Hash() *felt.Felt {
	return crypto.PedersenArray(
		c.APIVersion,
		crypto.PedersenArray(flatten(c.Externals)...),
		crypto.PedersenArray(flatten(c.L1Handlers)...),
		crypto.PedersenArray(flatten(c.Constructors)...),
		crypto.PedersenArray(c.Builtins...),
		c.ProgramHash,
		crypto.PedersenArray(c.Bytecode...),
	)
}

func flatten(entryPoints []EntryPoint) []*felt.Felt {
	result := make([]*felt.Felt, len(entryPoints)*2)
	for i, entryPoint := range entryPoints {
		// It is important that Selector is first because it
		// influences the class hash.
		result[2*i] = entryPoint.Selector
		result[2*i+1] = entryPoint.Offset
	}
	return result
}

// EntryPoint uniquely identifies a Cairo function to execute.
type EntryPoint struct {
	// starknet_keccak hash of the function signature.
	Selector *felt.Felt
	// The offset of the instruction in the class's bytecode.
	Offset *felt.Felt
}

type (
	Hints       map[uint64]interface{}
	Identifiers map[string]struct {
		CairoType   string         `json:"cairo_type,omitempty"`
		Decorators  *[]interface{} `json:"decorators,omitempty"`
		Destination string         `json:"destination,omitempty"`
		FullName    string         `json:"full_name,omitempty"`
		Members     *interface{}   `json:"members,omitempty"`
		Pc          *uint64        `json:"pc,omitempty"`
		References  *[]interface{} `json:"references,omitempty"`
		Size        *uint64        `json:"size,omitempty"`
		Type        string         `json:"type,omitempty"`
		Value       json.Number    `json:"value,omitempty"`
	}
	Program struct {
		Attributes       interface{} `json:"attributes,omitempty"`
		Builtins         []string    `json:"builtins"`
		CompilerVersion  string      `json:"compiler_version,omitempty"`
		Data             []string    `json:"data"`
		DebugInfo        interface{} `json:"debug_info"`
		Hints            Hints       `json:"hints"`
		Identifiers      Identifiers `json:"identifiers"`
		MainScope        interface{} `json:"main_scope"`
		Prime            string      `json:"prime"`
		ReferenceManager interface{} `json:"reference_manager"`
	}
)

type ContractDefinition struct {
	Abi         interface{} `json:"abi"`
	EntryPoints struct {
		Constructor []EntryPoint `json:"CONSTRUCTOR"`
		External    []EntryPoint `json:"EXTERNAL"`
		L1Handler   []EntryPoint `json:"L1_HANDLER"`
	} `json:"entry_points_by_type"`
	Program Program `json:"program"`
}

type ContractCode struct {
	Abi     interface{} `json:"abi"`
	Program Program     `json:"program"`
}

func ProgramHash(contractDefinition []byte) (*felt.Felt, error) {
	definition := new(ContractDefinition)
	err := json.Unmarshal(contractDefinition, &definition)
	if err != nil {
		return nil, err
	}

	program := definition.Program

	// make debug info None
	program.DebugInfo = nil

	// Cairo 0.8 added "accessible_scopes" and "flow_tracking_data" attribute fields, which were
	// not present in older contracts. They present as null/empty for older contracts deployed
	// prior to adding this feature and should not be included in the hash calculation in these cases.
	//
	// We therefore check and remove them from the definition before calculating the hash.
	if program.Attributes != nil {
		attributes := program.Attributes.([]interface{})
		if len(attributes) == 0 {
			program.Attributes = nil
		} else {
			for key, attribute := range attributes {
				attributeInterface := attribute.(map[string]interface{})
				if attributeInterface["accessible_scopes"] == nil || len(attributeInterface["accessible_scopes"].([]interface{})) == 0 {
					delete(attributeInterface, "accessible_scopes")
				}

				if attributeInterface["flow_tracking_data"] == nil || len(attributeInterface["flow_tracking_data"].(map[string]interface{})) == 0 {
					delete(attributeInterface, "flow_tracking_data")
				}

				attributes[key] = attributeInterface
			}

			program.Attributes = attributes
		}
	}

	contractCode := new(ContractCode)
	contractCode.Abi = definition.Abi
	contractCode.Program = program

	// Convert update program to bytes
	programBytes, err := contractCode.MarshalsJSON(definition.Program.Identifiers)
	if err != nil {
		return nil, err
	}

	programKeccak, err := crypto.StarkNetKeccak(programBytes)
	if err != nil {
		return nil, err
	}

	return programKeccak, nil
}

// MarshalsJSON is a custom json marshaller
func (c ContractCode) MarshalsJSON(identifiers Identifiers) ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.Write([]byte("{"))
	buf.Write([]byte("\"abi\": "))
	err := formatter(buf, c.Abi, "", identifiers)
	if err != nil {
		return nil, err
	}

	buf.Write([]byte(", "))
	buf.Write([]byte("\"program\": " + "{"))
	program := c.Program
	if program.Attributes != nil {
		buf.Write([]byte("\"attributes\": "))
		formatter(buf, program.Attributes, "", identifiers)
		buf.Write([]byte(", "))
	}
	buf.Write([]byte("\"builtins\": "))
	formatter(buf, program.Builtins, "", identifiers)
	buf.Write([]byte(", "))
	if program.CompilerVersion != "" {
		buf.Write([]byte("\"compiler_version\": "))
		formatter(buf, program.CompilerVersion, "", identifiers)
		buf.Write([]byte(", "))
	}
	buf.Write([]byte("\"data\": "))
	formatter(buf, program.Data, "", identifiers)
	buf.Write([]byte(", "))
	buf.Write([]byte("\"debug_info\": "))
	formatter(buf, program.DebugInfo, "", identifiers)
	buf.Write([]byte(", "))
	buf.Write([]byte("\"hints\": "))
	formatter(buf, program.Hints, "", identifiers)
	buf.Write([]byte(", "))
	buf.Write([]byte("\"identifiers\": "))

	identifiersInterface := make(map[string]interface{})
	for key, value := range identifiers {
		newValue, _ := json.Marshal(&value)
		var newValueInterface map[string]interface{}
		err := json.Unmarshal(newValue, &newValueInterface)
		if err != nil {
			return nil, err
		}
		identifiersInterface[key] = newValueInterface
	}
	formatter(buf, identifiersInterface, "", identifiers)
	buf.Write([]byte(", "))
	buf.Write([]byte("\"main_scope\": "))
	formatter(buf, program.MainScope, "", identifiers)
	buf.Write([]byte(", "))
	buf.Write([]byte("\"prime\": "))
	formatter(buf, program.Prime, "", identifiers)
	buf.Write([]byte(", "))
	buf.Write([]byte("\"reference_manager\": "))
	formatter(buf, program.ReferenceManager, "", identifiers)
	buf.Write([]byte("}}"))

	return buf.Bytes(), nil
}

// Explanation of why we have the tree and identifier parameters
//
// During marshalling, json converts big numbers (e.g. -106710729501573572985208420194530329073740042555888586719489)
// into floats (e.g. -1.0671072950157357e+59), which affects the computation.
//
// The extra extra parameters helps us to trace where these conversions happen
// and to avoid using floats in the hash computation.
func formatter(buf *bytes.Buffer, value interface{}, tree string, identifiers Identifiers) error {
	switch v := value.(type) {
	case string:
		var result string
		if strings.ContainsAny(v, "\n\\") {
			result = strings.ReplaceAll(v, "\\", "\\\\")
			result = strings.ReplaceAll(result, "\n", "\\n")
		} else {
			result = v
		}
		buf.WriteString("\"" + result + "\"")
	case uint:
		result := strconv.FormatUint(uint64(v), 10)
		buf.WriteString(result)
	case uint64:
		result := strconv.FormatUint(v, 10)
		buf.WriteString(result)
	case float64:
		result := strconv.FormatFloat(v, 'f', 0, 64)
		buf.WriteString(result)
	case json.Number:
		result := string(v)
		buf.WriteString(result)
	case map[string]interface{}:
		buf.Write([]byte{'{'})
		// Arrange lexicographically
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, k := range keys {
			buf.Write([]byte{'"'})
			buf.WriteString(k)
			buf.Write([]byte{'"'})
			buf.Write([]byte(": "))
			if k == "value" && reflect.TypeOf(v[k]) == reflect.TypeOf(float64(0)) && v[k] != float64(0) {
				value := identifiers[tree].Value
				buf.Write([]byte(value))
			} else {
				if reflect.TypeOf(v[k]) == reflect.TypeOf(map[string]interface{}{}) {
					tree = k
				}
				err := formatter(buf, v[k], tree, identifiers)
				if err != nil {
					return err
				}
			}
			if i < len(keys)-1 {
				buf.Write([]byte(", "))
			}
		}
		buf.Write([]byte{'}'})
	case Hints:
		buf.Write([]byte{'{'})
		// Arrange numerically
		keys := make([]uint64, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

		for i, k := range keys {
			buf.Write([]byte{'"'})
			buf.WriteString(strconv.FormatUint(k, 10))
			buf.Write([]byte{'"'})
			buf.Write([]byte(": "))
			err := formatter(buf, v[k], "", identifiers)
			if err != nil {
				return err
			}
			if i < len(keys)-1 {
				buf.Write([]byte(", "))
			}
		}
		buf.Write([]byte{'}'})
	case []interface{}:
		buf.Write([]byte{'['})
		count := 0
		for _, value := range v {
			err := formatter(buf, value, "", identifiers)
			if err != nil {
				return err
			}
			if count < len(v)-1 {
				buf.Write([]byte(", "))
			}
			count++
		}
		buf.Write([]byte{']'})
	case []string:
		buf.Write([]byte{'['})
		count := 0
		for _, value := range v {
			buf.WriteString(`"` + value + `"`)
			if count < len(v)-1 {
				buf.Write([]byte(", "))
			}
			count++
		}
		buf.Write([]byte{']'})
	default:
		if value == nil {
			buf.WriteString("null")
		} else {
			return fmt.Errorf("unknown type: %T", value)
		}
	}
	return nil
}

// Contract is an instance of a [Class].
type Contract struct {
	// The number of transactions sent from this contract.
	// Only account contracts can have a non-zero nonce.
	Nonce uint
	// Hash of the class that this contract instantiates.
	ClassHash *felt.Felt
	// Root of the contract's storage trie.
	StorageRoot *felt.Felt // TODO: is this field necessary?
}
