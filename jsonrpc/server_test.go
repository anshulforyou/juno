package jsonrpc

import (
	"errors"
	"reflect"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	failingTests := map[string]struct {
		data []byte
		want string
	}{
		"empty req":    {[]byte(""), "EOF"},
		"empty object": {[]byte("{}"), "unsupported RPC request version"},
		"wrong version": {[]byte(`
				{
					"jsonrpc" : "1.0"
				}`), "unsupported RPC request version"},
		"no method": {[]byte(`
				{
					"jsonrpc" : "2.0"
				}`), "no method specified"},
		"number param": {[]byte(`
				{
					"jsonrpc" : "2.0",
					"method" : "rpc_call",
					"params" : 44
				}`), "params should be an array or an object"},
		"string param": {[]byte(`
				{
					"jsonrpc" : "2.0",
					"method" : "rpc_call",
					"params" : "44"
				}`), "params should be an array or an object"},
		"array id": {[]byte(`
				{
					"jsonrpc" : "2.0",
					"method" : "rpc_call",
					"params" : { "malatya" : "44"},
					"id"     : [37]
				}`), "id should be a string or an integer"},
		"map id": {[]byte(`
				{
					"jsonrpc" : "2.0",
					"method" : "rpc_call",
					"params" : { "malatya" : "44"},
					"id"     : { "44" : "37"}
				}`), "id should be a string or an integer"},
		"float id": {[]byte(`
				{
					"jsonrpc" : "2.0",
					"method" : "rpc_call",
					"params" : { "malatya" : "44"},
					"id"     : 44.37
				}`), "id should be a string or an integer"},
	}

	for desc, test := range failingTests {
		_, err := newRequest(test.data)
		assert.EqualError(t, err, test.want, desc)
	}

	happyPathTests := map[string]struct {
		data []byte
	}{
		"uint id": {[]byte(`
				{
					"jsonrpc" : "2.0",
					"method" : "rpc_call",
					"params" : { "malatya" : "44"},
					"id"     : 44
				}`)},
		"int id": {[]byte(`
				{
					"jsonrpc" : "2.0",
					"method" : "rpc_call",
					"params" : { "malatya" : "44"},
					"id"     : -44
				}`)},
		"string id": {[]byte(`
				{
					"jsonrpc" : "2.0",
					"method" : "rpc_call",
					"params" : { "malatya" : 44, "depdep" : { "plate" : "37" }},
					"id"     : "-44"
				}`)},
	}

	for desc, test := range happyPathTests {
		_, err := newRequest(test.data)
		assert.NoError(t, err, desc)
	}
}

func TestServer_RegisterMethod(t *testing.T) {
	server := NewServer()
	tests := map[string]struct {
		handler    any
		paramNames []string
		want       string
	}{
		"not a func handler": {
			handler: 44,
			want:    "handler must be a function",
		},
		"excess param names": {
			handler:    func() {},
			paramNames: []string{"param1"},
			want:       "number of function params and param names must match",
		},
		"missing param names": {
			handler:    func(param1, param2 int) {},
			paramNames: []string{"param1"},
			want:       "number of function params and param names must match",
		},
		"no return": {
			handler: func(param1, param2 int) {},
			want:    "handler must return 2 values",
		},
		"int return": {
			handler: func(param1, param2 int) (int, int) { return 0, 0 },
			want:    "first return value must be an interface",
		},
		"no error return": {
			handler: func(param1, param2 int) (any, int) { return 0, 0 },
			want:    "second return value must be an error",
		},
	}

	for desc, test := range tests {
		err := server.RegisterMethod("method", test.paramNames, test.handler)
		assert.EqualError(t, err, test.want, desc)
	}

	err := server.RegisterMethod("method", nil, func(param1, param2 int) (any, error) { return 0, nil })
	assert.NoError(t, err)
}

func TestBuildArguments(t *testing.T) {
	type arg struct {
		Val1 int        `json:"val_1"`
		Val2 string     `json:"val_2"`
		Val3 *felt.Felt `json:"val_3"`
	}

	paramNames := []string{"param1", "param2"}
	felt, _ := new(felt.Felt).SetString("0x4437")

	tests := []struct {
		handler    any
		req        string
		shouldFail bool
		errorMsg   string
	}{
		{
			handler: func(param1 int, param2 *arg) (any, error) {
				assert.Equal(t, 44, param1)
				assert.Equal(t, true, param2 == nil)
				return 0, nil
			},
			req: `{
					"jsonrpc" : "2.0",
					"method" : "method",
					"params" : { 
						"param1" : 44 
					},
					"id"     : 44
				}`,
		},
		{
			handler: func(param1 int, param2 *arg) (any, error) {
				assert.Equal(t, 44, param1)
				assert.Equal(t, arg{
					Val1: 37,
					Val2: "juno",
					Val3: felt,
				}, *param2)
				return 0, nil
			},
			req: `{
					"jsonrpc" : "2.0",
					"method" : "method",
					"params" : { 
						"param1" : 44, 
						"param2" : { 
							"val_1" : 37,
							"val_2" : "juno",
							"val_3" : "0x4437"
						}
					},
					"id"     : 44
				}`,
		},
		{
			handler: func(param1 int, param2 arg) (any, error) {
				assert.Equal(t, 44, param1)
				assert.Equal(t, arg{
					Val1: 37,
					Val2: "juno",
					Val3: felt,
				}, param2)
				return 0, nil
			},
			req: `{
					"jsonrpc" : "2.0",
					"method" : "method",
					"params" : { 
						"param1" : 44, 
						"param2" : { 
							"val_1" : 37,
							"val_2" : "juno",
							"val_3" : "0x4437"
						}
					},
					"id"     : 44
				}`,
		},
		{
			handler: func(param1 int, param2 arg) (any, error) {
				assert.Equal(t, 44, param1)
				assert.Equal(t, arg{
					Val1: 37,
					Val2: "juno",
					Val3: felt,
				}, param2)
				return 0, nil
			},
			req: `{
					"jsonrpc" : "2.0",
					"method" : "method",
					"params" : [ 44, { 
							"val_1" : 37,
							"val_2" : "juno",
							"val_3" : "0x4437"
						}],
					"id"     : 44
				}`,
		},
		{
			handler: func(param1 int, param2 arg) (any, error) {
				assert.Equal(t, true, false) // should never be called
				return 0, nil
			},
			req: `{
					"jsonrpc" : "2.0",
					"method" : "method",
					"params" : [ { 
							"val_1" : 37,
							"val_2" : "juno",
							"val_3" : "0x4437"
						}, 44],
					"id"     : 44
				}`,
			shouldFail: true,
			errorMsg:   "json: cannot unmarshal object into Go value of type int",
		},
		{
			handler: func(param1 int, param2 arg) (any, error) {
				assert.Equal(t, true, false) // should never be called
				return 0, nil
			},
			req: `{
					"jsonrpc" : "2.0",
					"method" : "method",
					"params" : [ { 
							"val_1" : 37,
							"val_2" : "juno",
							"val_3" : "0x4437"
						}],
					"id"     : 44
				}`,
			shouldFail: true,
			errorMsg:   "missing param in list",
		},
	}

	for _, test := range tests {
		req, err := newRequest([]byte(test.req))
		if err != nil {
			t.Error(err)
		}

		args, err := buildArguments(req.Params, test.handler, paramNames)
		if !test.shouldFail && err != nil {
			t.Error(err)
		} else if !test.shouldFail {
			reflect.ValueOf(test.handler).Call(args)
		} else {
			assert.EqualError(t, err, test.errorMsg)
		}
	}
}

func TestHandle(t *testing.T) {
	server := NewServer()
	err := server.RegisterMethod("method", []string{"num", "shouldError", "msg"},
		func(num *int, shouldError bool, msg string) (any, error) {
			if shouldError {
				return nil, errors.New(msg)
			}
			return struct {
				Doubled int `json:"doubled"`
			}{*num * 2}, nil
		})
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		req string
		res string
	}{
		{
			req: `{]`,
			res: `{"jsonrpc":"2.0","error":{"code":-32600,"message":"invalid character ']' looking for beginning of object key string"},"id":null}`,
		},
		{
			req: `{"jsonrpc" : "1.0", "id" : 1}`,
			res: `{"jsonrpc":"2.0","error":{"code":-32600,"message":"unsupported RPC request version"},"id":null}`,
		},
		{
			req: `{"jsonrpc" : "2.0", "method" : "doesnotexits" , "id" : 2}`,
			res: `{"jsonrpc":"2.0","error":{"code":-32601,"message":"method not found"},"id":2}`,
		},
		{
			req: `{"jsonrpc" : "2.0", "method" : "method", "params" : [3, false] , "id" : 3}`,
			res: `{"jsonrpc":"2.0","error":{"code":-32602,"message":"missing param in list"},"id":3}`,
		},
		{
			req: `{"jsonrpc" : "2.0", "method" : "method", "params" : [3, false, "error message"] , "id" : 3}`,
			res: `{"jsonrpc":"2.0","result":{"doubled":6},"id":3}`,
		},
		{
			req: `{"jsonrpc" : "2.0", "method" : "method", "params" : [3, true, "error message"] , "id" : 4}`,
			res: `{"jsonrpc":"2.0","error":{"code":-32603,"message":"error message"},"id":4}`,
		},
		{
			req: `{"jsonrpc" : "2.0", "method" : "method",
					"params" : { "num" : 5, "shouldError" : false, "msg": "error message" } , "id" : 5}`,
			res: `{"jsonrpc":"2.0","result":{"doubled":10},"id":5}`,
		},
		{
			req: `{"jsonrpc" : "2.0", "method" : "method",
					"params" : { "num" : 5 } , "id" : 5}`,
			res: `{"jsonrpc":"2.0","result":{"doubled":10},"id":5}`,
		},
		{
			req: `{"jsonrpc" : "2.0", "method" : "method",
					"params" : { "num" : 5, "shouldError" : true } , "id" : 22}`,
			res: `{"jsonrpc":"2.0","error":{"code":-32603,"message":""},"id":22}`,
		},
	}

	for _, test := range tests {
		res, err := server.Handle([]byte(test.req))
		assert.NoError(t, err)
		assert.Equal(t, test.res, string(res))
	}
}
