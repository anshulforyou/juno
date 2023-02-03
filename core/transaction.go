package core

import (
	"errors"

	"github.com/NethermindEth/juno/core/crypto"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/core/trie"
	"github.com/NethermindEth/juno/utils"
	"github.com/ethereum/go-ethereum/common"
)

type Event struct {
	Data []*felt.Felt
	From *felt.Felt
	Keys []*felt.Felt
}

type L1ToL2Message struct {
	From     common.Address
	Nonce    *felt.Felt
	Payload  []*felt.Felt
	Selector *felt.Felt
	To       *felt.Felt
}

type L2ToL1Message struct {
	From    *felt.Felt
	Payload []*felt.Felt
	To      common.Address
}

type ExecutionResources struct {
	BuiltinInstanceCounter BuiltinInstanceCounter
	MemoryHoles            uint64
	Steps                  uint64
}

type BuiltinInstanceCounter struct {
	Bitwise    uint64
	EcOp       uint64
	Ecsda      uint64
	Output     uint64
	Pedersen   uint64
	RangeCheck uint64
}

type TransactionReceipt struct {
	ActualFee          *felt.Felt
	Events             []*Event
	ExecutionResources *ExecutionResources
	L1ToL2Message      *L1ToL2Message
	L2ToL1Message      []*L2ToL1Message
	TransactionHash    *felt.Felt
}

type Transaction interface {
	hash() *felt.Felt
	signatures() []*felt.Felt
	Type() string
}

type DeployTransaction struct {
	Hash *felt.Felt
	// A random number used to distinguish between different instances of the contract.
	ContractAddressSalt *felt.Felt
	// The address of the contract.
	ContractAddress *felt.Felt
	// The hash of the class which defines the contract’s functionality.
	ClassHash *felt.Felt
	// The arguments passed to the constructor during deployment.
	ConstructorCallData []*felt.Felt
	// The transaction’s version. Possible values are 1 or 0.
	//
	// When the fields that comprise a transaction change,
	// either with the addition of a new field or the removal of an existing field,
	// then the transaction version increases.
	// Transaction version 0 is deprecated and will be removed in a future version of StarkNet.
	Version *felt.Felt
}

func (d *DeployTransaction) Type() string {
	return "DEPLOY"
}

func (d *DeployTransaction) hash() *felt.Felt {
	return d.Hash
}

func (d *DeployTransaction) signatures() []*felt.Felt {
	return make([]*felt.Felt, 0)
}

type InvokeTransaction struct {
	Hash *felt.Felt
	// Version 0 fields
	// The address of the contract invoked by this transaction.
	ContractAddress *felt.Felt
	// The encoding of the selector for the function invoked (the entry point in the contract)
	EntryPointSelector *felt.Felt

	// Version 1 fields
	// The address of the sender of this transaction.
	SenderAddress *felt.Felt
	// The transaction nonce.
	Nonce *felt.Felt
	// The arguments that are passed to the validated and execute functions.
	CallData []*felt.Felt
	// Additional information given by the sender, used to validate the transaction.
	Signature []*felt.Felt
	// The maximum fee that the sender is willing to pay for the transaction
	MaxFee *felt.Felt
	// When the fields that comprise a transaction change,
	// either with the addition of a new field or the removal of an existing field,
	// then the transaction version increases.
	Version *felt.Felt
}

func (i *InvokeTransaction) Type() string {
	return "INVOKE_FUNCTION"
}

func (i *InvokeTransaction) hash() *felt.Felt {
	return i.Hash
}

func (i *InvokeTransaction) signatures() []*felt.Felt {
	return i.Signatures
}

type DeclareTransaction struct {
	Hash *felt.Felt
	// The class hash
	ClassHash *felt.Felt
	// The address of the account initiating the transaction.
	SenderAddress *felt.Felt
	// The maximum fee that the sender is willing to pay for the transaction.
	MaxFee *felt.Felt
	// Additional information given by the sender, used to validate the transaction.
	Signature []*felt.Felt
	// The transaction nonce.
	Nonce *felt.Felt
	// The transaction’s version. Possible values are 1 or 0.
	// When the fields that comprise a transaction change,
	// either with the addition of a new field or the removal of an existing field,
	// then the transaction version increases.
	// Transaction version 0 is deprecated and will be removed in a future version of StarkNet.
	Version *felt.Felt
}

func (d *DeclareTransaction) Type() string {
	return "DECLARE"
}

func (d *DeclareTransaction) hash() *felt.Felt {
	return d.Hash
}

func (d *DeclareTransaction) signatures() []*felt.Felt {
	return d.Signatures
}

func TransactionHash(transaction Transaction, network utils.Network) (*felt.Felt, error) {
	switch t := transaction.(type) {
	case *DeclareTransaction:
		return declareTransactionHash(t, network)
	case *InvokeTransaction:
		return invokeTransactionHash(t, network)
	case *DeployTransaction:
		return deployTransactionHash(t, network)
	default:
		return nil, errors.New("unknown transaction type")
	}
}

func deployTransactionHash(d *DeployTransaction, network utils.Network) (*felt.Felt, error) {
	snKeccakConstructor, err := crypto.StarkNetKeccak([]byte("constructor"))
	if err != nil {
		return nil, err
	}
	return crypto.PedersenArray(
		new(felt.Felt).SetBytes([]byte("deploy")),
		d.Version,
		d.ContractAddress,
		snKeccakConstructor,
		crypto.PedersenArray(d.ConstructorCallData...),
		new(felt.Felt),
		network.ChainId(),
	), nil
}

func invokeTransactionHash(i *InvokeTransaction, network utils.Network) (*felt.Felt, error) {
	invokeFelt := new(felt.Felt).SetBytes([]byte("invoke"))
	if i.Version.IsZero() {
		return crypto.PedersenArray(
			invokeFelt,
			i.ContractAddress,
			i.EntryPointSelector,
			crypto.PedersenArray(i.CallData...),
			network.ChainId(),
		), nil
	} else if i.Version.IsOne() {
		return crypto.PedersenArray(
			invokeFelt,
			i.Version,
			i.SenderAddress,
			new(felt.Felt),
			crypto.PedersenArray(i.CallData...),
			i.MaxFee,
			network.ChainId(),
			i.Nonce,
		), nil
	}
	return nil, errors.New("invalid transaction version")
}

func declareTransactionHash(d *DeclareTransaction, network utils.Network) (*felt.Felt, error) {
	declareFelt := new(felt.Felt).SetBytes([]byte("declare"))
	if d.Version.IsZero() {
		return crypto.PedersenArray(
			declareFelt,
			d.Version,
			d.SenderAddress,
			new(felt.Felt),
			crypto.PedersenArray(make([]*felt.Felt, 0)...),
			d.MaxFee,
			network.ChainId(),
			d.ClassHash,
		), nil
	} else if d.Version.IsOne() {
		return crypto.PedersenArray(
			declareFelt,
			d.Version,
			d.SenderAddress,
			new(felt.Felt),
			crypto.PedersenArray(d.ClassHash),
			d.MaxFee,
			network.ChainId(),
			d.Nonce,
		), nil
	}
	return nil, errors.New("invalid transaction version")
}

const commitmentTrieHeight uint = 64

// TransactionCommitment is the root of a height 64 binary Merkle Patricia tree of the
// transaction hashes and signatures in a block.
func TransactionCommitment(transactions []Transaction) (*felt.Felt, error) {
	var transactionCommitment *felt.Felt
	return transactionCommitment, trie.RunOnTempTrie(commitmentTrieHeight, func(trie *trie.Trie) error {
		zeroFelt := new(felt.Felt)
		for i, transaction := range transactions {
			signatures := transaction.signatures()
			hash := transaction.hash()

			signaturesHash := crypto.Pedersen(zeroFelt, zeroFelt)
			if len(signatures) > 0 {
				signaturesHash = crypto.PedersenArray(signatures...)
			}

			if _, err := trie.Put(new(felt.Felt).SetUint64(uint64(i)), crypto.Pedersen(hash, signaturesHash)); err != nil {
				return err
			}
		}
		root, err := trie.Root()
		if err != nil {
			return err
		}
		transactionCommitment = root
		return nil
	})
}

// EventCommitmentAndCount computes the event commitment and event count for a block.
func EventCommitmentAndCount(receipts []*TransactionReceipt) (*felt.Felt, uint64, error) {
	var eventCommitment *felt.Felt // root of a height 64 binary Merkle Patricia tree of the events in a block.
	var eventCount uint64          // number of events in a block.
	return eventCommitment, eventCount, trie.RunOnTempTrie(commitmentTrieHeight, func(trie *trie.Trie) error {
		for _, receipt := range receipts {
			for _, event := range receipt.Events {
				eventHash := crypto.PedersenArray(
					event.From,
					crypto.PedersenArray(event.Keys...),
					crypto.PedersenArray(event.Data...),
				)

				if _, err := trie.Put(new(felt.Felt).SetUint64(eventCount), eventHash); err != nil {
					return err
				}
				eventCount++
			}
		}
		root, err := trie.Root()
		if err != nil {
			return err
		}
		eventCommitment = root
		return nil
	})
}
