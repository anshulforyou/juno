package datasource

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/NethermindEth/juno/clients"
	"github.com/NethermindEth/juno/core/crypto"
	"github.com/NethermindEth/juno/core/felt"
)

type ContractCode struct {
	Abi     interface{}
	Program *clients.Program
}

func ProgramHash(contractDefinition *clients.ClassDefinition) (*felt.Felt, error) {
	program := contractDefinition.Program

	// make debug info None
	program.DebugInfo = nil

	// Cairo 0.8 added "accessible_scopes" and "flow_tracking_data" attribute fields, which were
	// not present in older contracts. They present as null/empty for older contracts deployed
	// prior to adding this feature and should not be included in the hash calculation in these cases.
	//
	// We, therefore, check and remove them from the definition before calculating the hash.
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
	contractCode.Abi = contractDefinition.Abi
	contractCode.Program = &program

	programBytes, err := contractCode.MarshalsJSON(contractDefinition.Program.Identifiers)
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
func (c ContractCode) MarshalsJSON(identifiers clients.Identifiers) ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.Write([]byte("{"))
	buf.Write([]byte("\"abi\": "))
	if err := formatter(buf, c.Abi, "", identifiers); err != nil {
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
		if err := json.Unmarshal(newValue, &newValueInterface); err != nil {
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
// The extra parameters help us to trace where these conversions happen
// and to avoid using floats in the hash computation.
func formatter(buf *bytes.Buffer, value interface{}, tree string, identifiers clients.Identifiers) error {
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
				if err := formatter(buf, v[k], tree, identifiers); err != nil {
					return err
				}
			}
			if i < len(keys)-1 {
				buf.Write([]byte(", "))
			}
		}
		buf.Write([]byte{'}'})
	case clients.Hints:
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
			if err := formatter(buf, v[k], "", identifiers); err != nil {
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
			if err := formatter(buf, value, "", identifiers); err != nil {
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
