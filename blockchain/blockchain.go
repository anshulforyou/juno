package blockchain

import (
	"sync"

	"github.com/NethermindEth/juno/core"
)

// Blockchain is responsible for keeping track of all things related to the StarkNet blockchain
type Blockchain struct {
	sync.RWMutex

	height *uint64
	// todo: much more
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		RWMutex: sync.RWMutex{},
	}
}

// Height returns the latest block height
func (b *Blockchain) Height() uint64 {
	b.RLock()
	defer b.RUnlock()
	return *b.height
}

// NextHeight returns the current height plus 1
func (b *Blockchain) NextHeight() uint64 {
	b.RLock()
	defer b.RUnlock()

	if b.height == nil {
		return 0
	}
	return *b.height + 1
}

func (b *Blockchain) Verify(block *core.Block) error {
	return nil
}

func (b *Blockchain) UpdateState(su *core.StateUpdate) error {
	return nil
}

func (b *Blockchain) Store(block *core.Block, su *core.StateUpdate) error {
	b.Lock()
	defer b.Unlock()
	*b.height++
	return nil
}
