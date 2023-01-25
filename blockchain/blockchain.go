package blockchain

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync"

	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/core/state"
	"github.com/NethermindEth/juno/db"
	"github.com/NethermindEth/juno/utils"
	"github.com/fxamacker/cbor/v2"
)

// Blockchain is responsible for keeping track of all things related to the StarkNet blockchain
type Blockchain struct {
	sync.RWMutex

	head     *core.Block
	network  utils.Network
	database db.DB
	// todo: much more
}

func NewBlockchain(database db.DB, network utils.Network) *Blockchain {
	// Todo: get the latest block from db using the prefix created in db/buckets.go
	return &Blockchain{
		RWMutex:  sync.RWMutex{},
		database: database,
		network:  network,
	}
}

// Height returns the latest block height
func (b *Blockchain) Height() uint64 {
	b.RLock()
	defer b.RUnlock()
	return b.head.Number
}

// NextHeight returns the current height plus 1
func (b *Blockchain) NextHeight() uint64 {
	b.RLock()
	defer b.RUnlock()

	if b.head == nil {
		return 0
	}
	return b.head.Number + 1
}

type blockDbKey struct {
	Number uint64
	Hash   *felt.Felt
}

func (k *blockDbKey) MarshalBinary() ([]byte, error) {
	var numB [8]byte
	binary.BigEndian.PutUint64(numB[:], k.Number)
	return db.Blocks.Key(numB[:], k.Hash.Marshal()), nil
}

func (k *blockDbKey) UnmarshalBinary(data []byte) error {
	if len(data) != 41 {
		return errors.New("key should be 41 bytes long")
	}

	if data[0] != byte(db.Blocks) {
		return errors.New("wrong prefix")
	}

	k.Number = binary.BigEndian.Uint64(data[1:9])
	k.Hash = new(felt.Felt).SetBytes(data[9:41])
	return nil
}

// Store takes a block and state update and performs sanity checks before putting in the database.
func (b *Blockchain) Store(block *core.Block, stateUpdate *core.StateUpdate) error {
	if err := b.verifyBlock(block, stateUpdate); err != nil {
		return err
	}
	if err := state.NewState(b.database).Update(stateUpdate); err != nil {
		return err
	}

	key := &blockDbKey{block.Number, block.Hash}
	bKey, err := key.MarshalBinary()
	if err != nil {
		return err
	}

	blockBinary, err := cbor.Marshal(block)
	if err != nil {
		return err
	}

	if err = b.database.Update(func(txn db.Transaction) error {
		if err = txn.Set(db.HeadBlock.Key(), blockBinary); err != nil {
			return err
		}
		return txn.Set(bKey, blockBinary)
	}); err != nil {
		return err
	}

	b.Lock()
	b.head = block
	b.Unlock()
	return nil
}

type ErrIncompatibleBlockAndStateUpdate struct {
	reason string
}

func (e ErrIncompatibleBlockAndStateUpdate) Error() string {
	return fmt.Sprintf("incompatible block and state update: %v", e.reason)
}

func (b *Blockchain) verifyBlock(block *core.Block, stateUpdate *core.StateUpdate) error {
	/*
		Todo: Transaction and TransactionReceipts
			- When Block is changed to include a list of Transaction and TransactionReceipts
			- Further checks would need to be added to ensure Transaction Hash has been computed
				properly.
			- Sanity check would need to include checks which ensure there is same number of
				Transactions and TransactionReceipts.
	*/
	b.RLock()
	if !block.ParentHash.Equal(b.head.Hash) {
		return errors.New("block's parent hash does not match head block hash")
	}
	b.RUnlock()
	if !block.Hash.Equal(stateUpdate.BlockHash) {
		return ErrIncompatibleBlockAndStateUpdate{"block hashes do not match"}
	}
	if !block.GlobalStateRoot.Equal(stateUpdate.NewRoot) {
		return ErrIncompatibleBlockAndStateUpdate{
			"block's GlobalStateRoot does not match state update's NewRoot",
		}
	}

	h, err := core.BlockHash(block, b.network)
	if !errors.Is(err, core.ErrUnverifiableBlock{}) {
		return err
	}

	if h != nil && !block.Hash.Equal(h) {
		return errors.New("incorrect block hash")
	}

	return nil
}
