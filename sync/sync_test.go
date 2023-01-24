package sync

import (
	"testing"

	"github.com/NethermindEth/juno/blockchain"
	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/stretchr/testify/assert"
)

/*

Todo:
- Get the fake StarkNetData to serve one block
	- For the above you need to have a head of the block in blockchain
	- You need to make the blockchain accept a db/transaction
	- Then get it to store in the fake db
	- Write todos for the verify and Store function
	- Decide what to do for

Todo:
- Create a Gateway Client Interfacde
- Create a FakeGatewayClient for testing purposes
- Rename the Gateway to Feeder Gateway to make it more specific
- Make the StarkNetData accept GatewayClient Interface
	- Pass the Gateway To the Synchronizer

*/

func TestSyncBlocks(t *testing.T) {
	bc := blockchain.NewBlockchain()
	fakeData := &fakeStarkNetData{}
	synchronizer := NewSynchronizer(bc, fakeData)
	err := synchronizer.SyncBlocks()
	assert.NoError(t, err)
}

type fakeStarkNetData struct{}

func (f *fakeStarkNetData) BlockByNumber(blockNumber uint64) (*core.Block, error) {
	return nil, nil
}

func (f *fakeStarkNetData) Transaction(transactionHash *felt.Felt) (*core.Transaction, error) {
	return nil, nil
}

func (f *fakeStarkNetData) Class(classHash *felt.Felt) (*core.Class, error) {
	return nil, nil
}

func (f *fakeStarkNetData) StateUpdate(blockNumber uint64) (*core.StateUpdate, error) {
	return nil, nil
}
