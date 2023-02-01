package core

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/utils"
)

func TestDeployTransactions(t *testing.T) {
	tests := map[string]struct {
		input   DeployTransaction
		network utils.Network
	}{
		// https://alpha-mainnet.starknet.io/feeder_gateway/get_transaction?transactionHash=0x6486c6303dba2f364c684a2e9609211c5b8e417e767f37b527cda51e776e6f0
		"Deploy transaction": {
			input: DeployTransaction{
				Hash:                hexToFelt("0x6486c6303dba2f364c684a2e9609211c5b8e417e767f37b527cda51e776e6f0"),
				ContractAddress:     hexToFelt("0x3ec215c6c9028ff671b46a2a9814970ea23ed3c4bcc3838c6d1dcbf395263c3"),
				ContractAddressSalt: hexToFelt("0x74dc2fe193daf1abd8241b63329c1123214842b96ad7fd003d25512598a956b"),
				ClassHash:           hexToFelt("0x46f844ea1a3b3668f81d38b5c1bd55e816e0373802aefe732138628f0133486"),
				ConstructorCallData: [](*felt.Felt){
					hexToFelt("0x6d706cfbac9b8262d601c38251c5fbe0497c3a96cc91a92b08d91b61d9e70c4"),
					hexToFelt("0x79dc0da7c54b95f10aa182ad0a46400db63156920adb65eca2654c0945a463"),
					hexToFelt("0x2"),
					hexToFelt("0x6658165b4984816ab189568637bedec5aa0a18305909c7f5726e4a16e3afef6"),
					hexToFelt("0x6b648b36b074a91eee55730f5f5e075ec19c0a8f9ffb0903cefeee93b6ff328"),
				},
				CallerAddress: new(felt.Felt).SetUint64(0),
				Version:       new(felt.Felt).SetUint64(0),
			},
			network: utils.MAINNET,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			transactionHash, err := TransactionHash(&test.input, test.network)
			if err != nil {
				t.Errorf("no error expected but got %v", err)
			}
			if !transactionHash.Equal(test.input.Hash) {
				t.Errorf("wrong hash: got %s, want %s", transactionHash.Text(16),
					test.input.Hash.Text(16))
			}
		})
	}
}

func TestInvokeTransactions(t *testing.T) {
	tests := map[string]struct {
		input   InvokeTransaction
		network utils.Network
	}{
		// https://alpha-mainnet.starknet.io/feeder_gateway/get_transaction?transactionHash=0xf1d99fb97509e0dfc425ddc2a8c5398b74231658ca58b6f8da92f39cb739e
		"Invoke transaction version 0": {
			input: InvokeTransaction{
				Hash:               hexToFelt("0xf1d99fb97509e0dfc425ddc2a8c5398b74231658ca58b6f8da92f39cb739e"),
				ContractAddress:    hexToFelt("0x43324c97e376d7d164abded1af1e73e9ce8214249f711edb7059c1ca34560e8"),
				EntryPointSelector: hexToFelt("0x317eb442b72a9fae758d4fb26830ed0d9f31c8e7da4dbff4e8c59ea6a158e7f"),
				CallData: [](*felt.Felt){
					hexToFelt("0x1b654cb59f978da2eee76635158e5ff1399bf607cb2d05e3e3b4e41d7660ca2"),
					hexToFelt("0x2"),
					hexToFelt("0x5f743efdb29609bfc2002041bdd5c72257c0c6b5c268fc929a3e516c171c731"),
					hexToFelt("0x635afb0ea6c4cdddf93f42287b45b67acee4f08c6f6c53589e004e118491546"),
				},
				MaxFee:  hexToFelt("0x0"),
				Version: new(felt.Felt).SetUint64(0),
			},
			network: utils.MAINNET,
		},
		// https://alpha-mainnet.starknet.io/feeder_gateway/get_transaction?transactionHash=0x2897e3cec3e24e4d341df26b8cf1ab84ea1c01a051021836b36c6639145b497
		"Invoke transaction version 1": {
			input: InvokeTransaction{
				Hash:            hexToFelt("0x2897e3cec3e24e4d341df26b8cf1ab84ea1c01a051021836b36c6639145b497"),
				ContractAddress: hexToFelt("0x3ec215c6c9028ff671b46a2a9814970ea23ed3c4bcc3838c6d1dcbf395263c3"),
				CallData: [](*felt.Felt){
					hexToFelt("0x1"),
					hexToFelt("0x727a63f78ee3f1bd18f78009067411ab369c31dece1ae22e16f567906409905"),
					hexToFelt("0x22de356837ac200bca613c78bd1fcc962a97770c06625f0c8b3edeb6ae4aa59"),
					hexToFelt("0x0"),
					hexToFelt("0xb"),
					hexToFelt("0xb"),
					hexToFelt("0xa"),
					hexToFelt("0x6db793d93ce48bc75a5ab02e6a82aad67f01ce52b7b903090725dbc4000eaa2"),
					hexToFelt("0x6141eac4031dfb422080ed567fe008fb337b9be2561f479a377aa1de1d1b676"),
					hexToFelt("0x27eb1a21fa7593dd12e988c9dd32917a0dea7d77db7e89a809464c09cf951c0"),
					hexToFelt("0x400a29400a34d8f69425e1f4335e6a6c24ce1111db3954e4befe4f90ca18eb7"),
					hexToFelt("0x599e56821170a12cdcf88fb8714057ce364a8728f738853da61d5b3af08a390"),
					hexToFelt("0x46ad66f467df625f3b2dd9d3272e61713e8f74b68adac6718f7497d742cfb17"),
					hexToFelt("0x4f348b585e6c1919d524a4bfe6f97230ecb61736fe57534ec42b628f7020849"),
					hexToFelt("0x19ae40a095ffe79b0c9fc03df2de0d2ab20f59a2692ed98a8c1062dbf691572"),
					hexToFelt("0xe120336994adef6c6e47694f87278686511d4622997d4a6f216bd6e9fa9acc"),
					hexToFelt("0x56e6637a4958d062db8c8198e315772819f64d915e5c7a8d58a99fa90ff0742"),
				},
				Signature: [](*felt.Felt){
					hexToFelt("0x383ba105b6d0f59fab96a412ad267213ddcd899e046278bdba64cd583d680b"),
					hexToFelt("0x1896619a17fde468978b8d885ffd6f5c8f4ac1b188233b81b91bcf7dbc56fbd"),
				},
				Nonce:         hexToFelt("0x42"),
				SenderAddress: hexToFelt("0x1fc039de7d864580b57a575e8e6b7114f4d2a954d7d29f876b2eb3dd09394a0"),
				MaxFee:        hexToFelt("0x17f0de82f4be6"),
				Version:       new(felt.Felt).SetUint64(1),
			},
			network: utils.MAINNET,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			transactionHash, err := TransactionHash(&test.input, test.network)
			if err != nil {
				t.Errorf("no error expected but got %v", err)
			}
			if !transactionHash.Equal(test.input.Hash) {
				t.Errorf("wrong hash: got %s, want %s", transactionHash.Text(16),
					test.input.Hash.Text(16))
			}
		})
	}
}

func TestDeclareTransaction(t *testing.T) {
	tests := map[string]struct {
		input   DeclareTransaction
		network utils.Network
	}{
		// https://alpha-mainnet.starknet.io/feeder_gateway/get_transaction?transactionHash=0x222f8902d1eeea76fa2642a90e2411bfd71cffb299b3a299029e1937fab3fe4
		"Declare transaction version 0": {
			input: DeclareTransaction{
				Hash: hexToFelt("0x222f8902d1eeea76fa2642a90e2411bfd71cffb299b3a299029e1937fab3fe4"),
				// https://alpha-mainnet.starknet.io/feeder_gateway/get_class_by_hash?classHash=0x2760f25d5a4fb2bdde5f561fd0b44a3dee78c28903577d37d669939d97036a0
				ClassHash:     hexToFelt("0x2760f25d5a4fb2bdde5f561fd0b44a3dee78c28903577d37d669939d97036a0"),
				Nonce:         hexToFelt("0x0"),
				SenderAddress: hexToFelt("0x1"),
				MaxFee:        hexToFelt("0x0"),
				Version:       new(felt.Felt).SetUint64(0),
			},
			network: utils.MAINNET,
		},
		// https://alpha-mainnet.starknet.io/feeder_gateway/get_transaction?transactionHash=0x1b4d9f09276629d496af1af8ff00173c11ff146affacb1b5c858d7aa89001ae
		"Declare transaction version 1": {
			input: DeclareTransaction{
				Hash:      hexToFelt("0x1b4d9f09276629d496af1af8ff00173c11ff146affacb1b5c858d7aa89001ae"),
				ClassHash: hexToFelt("0x7aed6898458c4ed1d720d43e342381b25668ec7c3e8837f761051bf4d655e54"),
				Signature: [](*felt.Felt){
					hexToFelt("0x221b9576c4f7b46d900a331d89146dbb95a7b03d2eb86b4cdcf11331e4df7f2"),
					hexToFelt("0x667d8062f3574ba9b4965871eec1444f80dacfa7114e1d9c74662f5672c0620"),
				},
				Nonce:         hexToFelt("0x5"),
				SenderAddress: hexToFelt("0x39291faa79897de1fd6fb1a531d144daa1590d058358171b83eadb3ceafed8"),
				MaxFee:        hexToFelt("0xf6dbd653833"),
				Version:       new(felt.Felt).SetUint64(1),
			},
			network: utils.MAINNET,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			transactionHash, err := TransactionHash(&test.input, test.network)
			if err != nil {
				t.Errorf("no error expected but got %v", err)
			}
			if !transactionHash.Equal(test.input.Hash) {
				t.Errorf("wrong hash: got %s, want %s", transactionHash.Text(16),
					test.input.Hash.Text(16))
			}
		})
	}
}

var (
	//go:embed testdata/block_156000.json
	block156000 []byte
	//go:embed testdata/block_1.json
	block1Goerli []byte
	//go:embed testdata/block_1_integration.json
	block1Integration []byte
	//go:embed testdata/block_16789_main.json
	blocks16789Main []byte
)

var receipts [][]*TransactionReceipt

func getTransactionReceipts(t *testing.T) {
	// https://alpha4.starknet.io/feeder_gateway/get_block?blockNumber=156000
	var blck156000 map[string]interface{}
	if err := json.Unmarshal(block156000, &blck156000); err != nil {
		t.Fatal(err)
	}
	receiptsInterface := blck156000["transaction_receipts"].([]interface{})
	txns := blck156000["transactions"].([]interface{})
	receipt156000 := generateReceipt(txns, receiptsInterface)

	// https://alpha4.starknet.io/feeder_gateway/get_block?blockNumber=1
	var blck1Goerli map[string]interface{}
	if err := json.Unmarshal(block1Goerli, &blck1Goerli); err != nil {
		t.Fatal(err)
	}
	receiptsInterface = blck1Goerli["transaction_receipts"].([]interface{})
	txns = blck1Goerli["transactions"].([]interface{})
	receipt1Goerli := generateReceipt(txns, receiptsInterface)

	// https://external.integration.starknet.io/feeder_gateway/get_block?blockNumber=1
	var blck1Integration map[string]interface{}
	if err := json.Unmarshal(block1Integration, &blck1Integration); err != nil {
		t.Fatal(err)
	}
	receiptsInterface = blck1Integration["transaction_receipts"].([]interface{})
	txns = blck1Integration["transactions"].([]interface{})
	receipt1Integration := generateReceipt(txns, receiptsInterface)

	// https://alpha-mainnet.starknet.io/feeder_gateway/get_block?blockNumber=16789
	var blck16789Main map[string]interface{}
	if err := json.Unmarshal(blocks16789Main, &blck16789Main); err != nil {
		t.Fatal(err)
	}
	receiptsInterface = blck16789Main["transaction_receipts"].([]interface{})
	txns = blck16789Main["transactions"].([]interface{})
	receipt16789Main := generateReceipt(txns, receiptsInterface)

	receipts = [][]*TransactionReceipt{
		receipt156000,
		receipt1Goerli,
		receipt1Integration,
		receipt16789Main,
	}
}

func generateReceipt(txns []interface{}, receiptsInterface []interface{}) []*TransactionReceipt {
	receipts := make([]*TransactionReceipt, len(txns))

	transactionType := func(t string) TransactionType {
		switch t {
		case "DECLARE":
			return Declare
		case "DEPLOY":
			return Deploy
		case "DEPLOY_ACCOUNT":
			return DeployAccount
		case "INVOKE_FUNCTION":
			return Invoke
		case "L1_HANDLER":
			return L1Handler
		default:
			return -1
		}
	}

	for i, r := range receiptsInterface {
		receipt := r.(map[string]interface{})
		txn := txns[i].(map[string]interface{})
		var events []*Event
		for _, e := range receipt["events"].([]interface{}) {
			event := e.(map[string]interface{})
			var data []*felt.Felt
			for _, d := range event["data"].([]interface{}) {
				data = append(data, hexToFelt(d.(string)))
			}
			var keys []*felt.Felt
			for _, k := range event["keys"].([]interface{}) {
				keys = append(keys, hexToFelt(k.(string)))
			}
			events = append(events, &Event{
				Data: data,
				From: hexToFelt(event["from_address"].(string)),
				Keys: keys,
			})
		}
		var signatures []*felt.Felt
		if txn["signature"] != nil {
			for _, s := range txn["signature"].([]interface{}) {
				signatures = append(signatures, hexToFelt(s.(string)))
			}
		}
		// Some of these values are set to nil since they are not required to calculate the commitment.
		transactionReceipt := TransactionReceipt{
			Events:          events,
			Signatures:      signatures,
			TransactionHash: hexToFelt(receipt["transaction_hash"].(string)),
			Type:            transactionType(txn["type"].(string)),
		}
		receipts[i] = &transactionReceipt
	}
	return receipts
}

func init() {
	var t *testing.T
	getTransactionReceipts(t)
}

func assertCorrectCommitment(t *testing.T, got *felt.Felt, want string) {
	t.Helper()
	if "0x"+got.Text(16) != want {
		t.Errorf("got %s, want %s", "0x"+got.Text(16), want)
	}
}

func TestTransactionCommitment(t *testing.T) {
	tests := []struct {
		description string
		receipts    []*TransactionReceipt
		want        string
	}{
		{
			"receipt 1 (goerli)",
			receipts[0],
			"0x24638e0ca122d0260d54e901dc0942ea68bd1fc40a96b5da765985c47c92500",
		},
		{
			"receipt 2 (goerli)",
			receipts[1],
			"0x18bb7d6c1c558aa0a025f08a7d723a44b13008ffb444c432077f319a7f4897c",
		},
		{
			"receipt 1 (integration)",
			receipts[2],
			"0xbf11745df434cbd284e13ca36354139a4bca2f6722e737c6136590990c8619",
		},
		{
			"receipt 1 (mainnet)",
			receipts[3],
			"0x580a06bfc8c3fe39bbb7c5d16298b8928bf7c28f4c31b8e6b48fc25cd644fc1",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			commitment, _ := TransactionCommitment(test.receipts)
			assertCorrectCommitment(t, commitment, test.want)
		})
	}
}

func TestEventCommitment(t *testing.T) {
	tests := []struct {
		description string
		receipts    []*TransactionReceipt
		want        string
	}{
		{
			"receipt 1 (goerli)",
			receipts[0],
			"0x5d25e41d43b00681cc63ed4e13a82efe3e02f47e03173efbd737dd52ba88c7e",
		},
		{
			"receipt 2 (goerli)",
			receipts[1],
			"0x0",
		},
		{
			"receipt 3 (integration)",
			receipts[2],
			"0x0",
		},
		{
			"receipt 4 (mainnet)",
			receipts[3],
			"0x6f499789aabb31935810ce89d6ea9e9d37c5921c0d7fae2bd68f2fff5b7b93f",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			commitment, _, _ := EventCommitmentAndCount(test.receipts)
			assertCorrectCommitment(t, commitment, test.want)
		})
	}
}
