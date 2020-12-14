/*
 *  Copyright 2018 KardiaChain
 *  This file is part of the go-kardia library.
 *
 *  The go-kardia library is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Lesser General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  The go-kardia library is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 *  GNU Lesser General Public License for more details.
 *
 *  You should have received a copy of the GNU Lesser General Public License
 *  along with the go-kardia library. If not, see <http://www.gnu.org/licenses/>.
 */
// Package kardia
package kardia

import (
	"context"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/blendle/zapdriver"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/kardiachain/go-kardiamain/lib/p2p"
	coreTypes "github.com/kardiachain/go-kardiamain/types"

	"github.com/kardiachain/explorer-backend/metrics"
	"github.com/kardiachain/explorer-backend/types"
)

type testSuite struct {
	cfg *Config

	minBlockNumber uint64

	m *metrics.Provider

	blockHeight uint64
	blockHash   string
	txHash      string
	address     string

	sampleBlock       *types.Block
	sampleBlockHeader *types.Header
	sampleTx          *types.Transaction
	sampleTxReceipt   *types.Receipt
	samplePeer        *p2p.PeerInfo
	sampleNodeInfo    *p2p.NodeInfo
	sampleDatadir     string
	sampleValidator   []*types.Validator
}

func setupTestSuite() *testSuite {
	blockHeight := uint64(395)
	blockHash := "0x634662e42bc71d2a7ca767ca19735c8f19694fd7dfbbc70bb28698e0e01be888"
	txHash := "0x02e90c26892a6d230b6964a78446e055b289c5ad53039468ea6a147c0ee31990"
	address := "0xc1fe56E3F58D3244F606306611a5d10c8333f1f6"
	sampleBlock := &types.Block{
		Hash:   blockHash,
		Height: blockHeight,
		Time:   1601908120,
		NumTxs: 0,
		// NumDualEvents:
		GasLimit:   1050000000,
		GasUsed:    0,
		LastBlock:  "0xf9fd47f388c3f41214d55c51fc3d59c8a5e550099a2aa3468d500539d15b0c7a",
		CommitHash: "0x45239add9675da72e0edb5c0ccc80d0f1c758a3cd398fee7f863d69d4863a759",
		DataHash:   "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		// DualEventsHash:
		ReceiptsRoot:  "",
		LogsBloom:     coreTypes.Bloom{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		Validator:     "0xc1fe56E3F58D3244F606306611a5d10c8333f1f6",
		ValidatorHash: "0x6231cec385931237749482972bf28d819fe9527c5ba618cd3620a1ba3be65bbd",
		ConsensusHash: "0x0000000000000000000000000000000000000000000000000000000000000000",
		AppHash:       "0x0d5d8f1a6fdffac4d9c93b1e230da611b65bad7b3c82fb28300247e3a40df76c",
		EvidenceHash:  "0x0000000000000000000000000000000000000000000000000000000000000000",

		Txs:      []*types.Transaction(nil),
		Receipts: []*types.Receipt(nil),
	}
	sampleBlockHeader := &types.Header{
		Hash:   blockHash,
		Height: blockHeight,
		Time:   1601908120,
		NumTxs: 0,
		// NumDualEvents:
		GasLimit:   1050000000,
		GasUsed:    0,
		LastBlock:  "0xf9fd47f388c3f41214d55c51fc3d59c8a5e550099a2aa3468d500539d15b0c7a",
		CommitHash: "0x45239add9675da72e0edb5c0ccc80d0f1c758a3cd398fee7f863d69d4863a759",
		DataHash:   "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		// DualEventsHash:
		ReceiptsRoot:  "",
		LogsBloom:     coreTypes.Bloom{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		Validator:     "0xc1fe56E3F58D3244F606306611a5d10c8333f1f6",
		ValidatorHash: "0x6231cec385931237749482972bf28d819fe9527c5ba618cd3620a1ba3be65bbd",
		ConsensusHash: "0x0000000000000000000000000000000000000000000000000000000000000000",
		AppHash:       "0x0d5d8f1a6fdffac4d9c93b1e230da611b65bad7b3c82fb28300247e3a40df76c",
		EvidenceHash:  "0x0000000000000000000000000000000000000000000000000000000000000000",
	}
	sampleTx := &types.Transaction{
		Hash: txHash,
		To:   "0x2500A193147c8B8FfB4808564a2DC0f476400B86",
		From: address,
		// Status:
		// ContractAddress:
		Value:    "2",
		GasPrice: 1,
		GasFee:   21000,
		// GasLimit:
		BlockNumber: 142,
		Nonce:       1304,
		BlockHash:   "0x484d030a20881754beea5b17485868df2e9cde3fea20adbe9ae48dbc73529605",
		Time:        1601908519,
		InputData:   "0x",
		// Logs:
		TransactionIndex: 1,
		// ReceiptReceived:
	}
	sampleTxReceipt := &types.Receipt{
		BlockHash:         "0x484d030a20881754beea5b17485868df2e9cde3fea20adbe9ae48dbc73529605",
		BlockHeight:       142,
		TransactionHash:   txHash,
		TransactionIndex:  1,
		From:              address,
		To:                "0x2500A193147c8B8FfB4808564a2DC0f476400B86",
		GasUsed:           21000,
		CumulativeGasUsed: 42000,
		ContractAddress:   "0x",
		Logs:              []types.Log{},
		LogsBloom:         coreTypes.Bloom{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		Root:              "0x",
		Status:            1,
	}
	samplePeer := &p2p.PeerInfo{}
	sampleNodeInfo := &p2p.NodeInfo{}
	sampleValidator := []*types.Validator{}
	m := metrics.New()
	cfg := &Config{
		rpcURL:            []string{"http://10.10.0.251:8545", "http://10.10.0.251:8546", "http://10.10.0.251:8547", "http://10.10.0.251:8548", "http://10.10.0.251:8549", "http://10.10.0.251:8550", "http://10.10.0.251:8551"},
		trustedNodeRPCURL: []string{"http://10.10.0.251:8545", "http://10.10.0.251:8546", "http://10.10.0.251:8547"},
	}
	loggerCfg := zapdriver.NewProductionConfig()
	logger, err := loggerCfg.Build()
	if err != nil {
		return nil
	}
	defer logger.Sync()
	cfg.lgr = logger
	return &testSuite{
		cfg:               cfg,
		minBlockNumber:    1<<bits.UintSize - 1,
		m:                 m,
		blockHeight:       blockHeight,
		blockHash:         blockHash,
		txHash:            txHash,
		address:           address,
		sampleBlock:       sampleBlock,
		sampleBlockHeader: sampleBlockHeader,
		sampleTx:          sampleTx,
		sampleTxReceipt:   sampleTxReceipt,
		samplePeer:        samplePeer,
		sampleNodeInfo:    sampleNodeInfo,
		sampleDatadir:     "/home/.kardia",
		sampleValidator:   sampleValidator,
	}
}

func SetupKAIClient() (ClientInterface, context.Context, *testSuite, error) {
	suite := setupTestSuite()
	ctx, cancelFn := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range sigCh {
			cancelFn()
		}
	}()
	client, err := NewKaiClient(suite.cfg)
	if err != nil {
		return nil, nil, suite, fmt.Errorf("Failed to create new KaiClient: %v", err)
	}
	return client, ctx, suite, nil
}

func TestSanity(t *testing.T) {

}

func TestLatestBlockNumber(t *testing.T) {
	client, ctx, _, err := SetupKAIClient()
	assert.Nil(t, err)

	num, err := client.LatestBlockNumber(ctx)

	assert.Nil(t, err)
	t.Log("Latest block number: ", num)
	assert.IsTypef(t, uint64(0), num, "Block number must be an uint64")
}

func TestBlockByHash(t *testing.T) {
	client, ctx, testSuite, err := SetupKAIClient()
	assert.Nil(t, err)

	b, err := client.BlockByHash(ctx, testSuite.blockHash)

	assert.Nil(t, err)
	t.Log("Hash: ", testSuite.blockHash, "\nBlock: ", b)
	assert.EqualValuesf(t, testSuite.sampleBlock, b, "Received block must be equal to sampleBlock in testSuite")
}

func TestBlockByNumber(t *testing.T) {
	client, ctx, testSuite, err := SetupKAIClient()
	assert.Nil(t, err)

	b, err := client.BlockByHeight(ctx, testSuite.blockHeight)

	assert.Nil(t, err)
	t.Log("\nBlock number: ", testSuite.blockHeight, "\nBlock: ", b)
	assert.EqualValuesf(t, testSuite.sampleBlock, b, "Received block must be equal to sampleBlock in testSuite")
}

func TestBlockHeaderByNumber(t *testing.T) {
	client, ctx, testSuite, err := SetupKAIClient()
	assert.Nil(t, err)

	h, err := client.BlockHeaderByNumber(ctx, testSuite.blockHeight)

	assert.Nil(t, err)
	t.Log("Block number: ", testSuite.blockHeight, "\nBlock header: ", h)
	assert.EqualValuesf(t, testSuite.sampleBlockHeader, h, "Received block header must be equal to sampleBlockHeader in testSuite")
}

func TestBlockHeaderByHash(t *testing.T) {
	client, ctx, testSuite, err := SetupKAIClient()
	assert.Nil(t, err)

	h, err := client.BlockHeaderByHash(ctx, testSuite.blockHash)

	assert.Nil(t, err)
	t.Log("\nHash: ", testSuite.blockHash, "\nBlock header: ", h)
	assert.EqualValuesf(t, testSuite.sampleBlockHeader, h, "Received block header must be equal to sampleBlockHeader in testSuite")
}

func TestBalanceAt(t *testing.T) {
	client, ctx, testSuite, err := SetupKAIClient()
	assert.Nil(t, err)

	b, err := client.BalanceAt(ctx, testSuite.address, nil)

	assert.Nil(t, err)
	t.Log("Address: ", testSuite.address, " Balance: ", b)
	assert.IsTypef(t, "", b, "Balance must be a string")
	assert.NotEqualValuesf(t, b, "-1", "Balance must be larger than -1")
}

func TestNonceAt(t *testing.T) {
	client, ctx, testSuite, err := SetupKAIClient()
	assert.Nil(t, err)

	n, err := client.NonceAt(ctx, testSuite.address)

	assert.Nil(t, err)
	t.Log("Address: ", testSuite.address, " Nonce: ", n)
	assert.IsTypef(t, uint64(0), n, "Nonce must be an uint64")
}

func TestGetTransaction(t *testing.T) {
	client, ctx, testSuite, err := SetupKAIClient()
	assert.Nil(t, err)

	tx, isPending, err := client.GetTransaction(ctx, testSuite.txHash)

	assert.Nil(t, err)
	assert.EqualValuesf(t, false, isPending, "isPending must be true")
	assert.EqualValuesf(t, testSuite.sampleTx, tx, "Received tx must be equal to sampleTx in testSuite")
}

func TestGetTransactionReceipt(t *testing.T) {
	client, ctx, testSuite, err := SetupKAIClient()
	assert.Nil(t, err)

	receipt, err := client.GetTransactionReceipt(ctx, testSuite.txHash)

	assert.Nil(t, err)
	assert.EqualValuesf(t, testSuite.sampleTxReceipt, receipt, "Received receipt must be equal to sampleTxReceipt in testSuite")
}

// createCustomClient create a Client instance which has a custom number of RPC nodes
func createCustomClient(numOfNodes int) (*Client, error) {
	testSuite := setupTestSuite()
	rpcUrls := []string{}
	for i := 0; i < numOfNodes; i++ {
		rpcUrls = append(rpcUrls, testSuite.cfg.rpcURL[i])
	}
	client, err := NewKaiClient(testSuite.cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new KaiClient: %v", err)
	}
	return client.(*Client), nil
}

// getTxsHashList get tx hashes of a random block to avoid data caching at RPC nodes
func getTxsHashList(ctx context.Context, client *Client, requiredNumOfTxs int) []string {
	txsHash := []string{}
	for {
		num, _ := client.LatestBlockNumber(ctx)
		blockHeight := uint64(rand.Intn(int(num)))
		b, _ := client.BlockByHeight(ctx, blockHeight)
		for _, tx := range b.Txs {
			txsHash = append(txsHash, tx.Hash)
			if len(txsHash) >= requiredNumOfTxs {
				break
			}
		}
		if len(txsHash) >= requiredNumOfTxs {
			break
		}
	}
	return txsHash
}

func workerGetTxReceipt(ctx context.Context, wg *sync.WaitGroup, client *Client, m *metrics.Provider, jobs <-chan string, logger *zap.Logger, workerIndex int) {
	defer wg.Done()
	startTime := time.Now()
	for j := range jobs {
		_, _ = client.GetTransactionReceipt(ctx, j)
	}
	endTime := time.Since(startTime)
	m.RecordInsertBlockTime(endTime)
	logger.Info("Result: ", zap.Int("Worker index: ", workerIndex), zap.String("Processing time: ", endTime.String()))
	return
}

// getTxReceipts benchmark a batch of get tx receipt calls with custom parameters
func getTxReceipts(m *metrics.Provider, numOfTx int, numOfWorker int, numOfNodes int) {
	ctx := context.Background()
	client, err := createCustomClient(numOfNodes)
	if err != nil {
		return
	}
	txsHashList := getTxsHashList(ctx, client, numOfTx)
	logger := client.lgr.With(zap.Int("numOfTx: ", len(txsHashList)), zap.Int("numOfWorker: ", numOfWorker), zap.Int("numOfNodes: ", numOfNodes))

	job := make(chan string, 10000)
	var wg sync.WaitGroup
	for i := 0; i < numOfWorker; i++ {
		wg.Add(1)
		go workerGetTxReceipt(ctx, &wg, client, m, job, logger, i)
	}

	for _, hash := range txsHashList {
		job <- hash
	}
	close(job)

	wg.Wait()
	logger.Info("Elapsed time: ", zap.String("Avg processing time: ", m.GetInsertBlockTime()))
	m.Reset()
}

func TestGetTransactionReceipt2(t *testing.T) {
	numOfTx := []int{1000, 10000}
	numOfWorker := []int{4, 8, 16}
	numOfNodes := []int{3, 5, 7} // mean 2, 4, 6 because one node will be excluded to retrive normal RPC requests
	m := metrics.New()

	for i := 0; i < len(numOfTx); i++ {
		for j := 0; j < len(numOfWorker); j++ {
			for k := 0; k < len(numOfNodes); k++ {
				getTxReceipts(m, numOfTx[i], numOfWorker[j], numOfNodes[k])
			}
		}
	}
}

func TestPeers(t *testing.T) {
	client, ctx, _, err := SetupKAIClient()
	assert.Nil(t, err)

	peers, err := client.Peers(ctx)

	assert.Nil(t, err)
	assert.IsTypef(t, []*p2p.PeerInfo{}, peers, "peers must be an array of *p2p.PeerInfo")
	// assert.EqualValuesf(t, testSuite.sampleTxReceipt, peers, "Received receipt must be equal to sampleTxReceipt in testSuite")
}

func TestNodesInfo(t *testing.T) {
	client, ctx, _, err := SetupKAIClient()
	assert.Nil(t, err)

	node, err := client.NodesInfo(ctx)

	assert.Nil(t, err)
	assert.IsTypef(t, &p2p.NodeInfo{}, node, "node must be a *p2p.NodeInfo")
	// assert.EqualValuesf(t, testSuite.sampleTxReceipt, node, "Received receipt must be equal to sampleTxReceipt in testSuite")
}

func TestDataDir(t *testing.T) {
	client, ctx, testSuite, err := SetupKAIClient()
	assert.Nil(t, err)

	dir, err := client.Datadir(ctx)

	assert.Nil(t, err)
	assert.EqualValuesf(t, testSuite.sampleDatadir, dir, "Receive data directory must be equal to sampleDatadir in testSuite")
}

func TestValidators(t *testing.T) {
	client, ctx, testSuite, err := SetupKAIClient()
	assert.Nil(t, err)

	validators, err := client.Validators(ctx)
	assert.Nil(t, err)

	assert.IsTypef(t, testSuite.sampleValidator, validators, "Validators must be a []*Validator")
	// assert.EqualValuesf(t, testSuite.sampleDatadir, dir, "Receive data directory must be equal to sampleDatadir in testSuite")
}

func TestValidator(t *testing.T) {
	client, ctx, _, err := SetupKAIClient()
	assert.Nil(t, err)

	validator, err := client.Validator(ctx, "")
	assert.Nil(t, err)

	assert.IsTypef(t, &types.Validator{}, validator, "Validator must be a *Validator")
	// assert.EqualValuesf(t, testSuite.sampleDatadir, dir, "Receive data directory must be equal to sampleDatadir in testSuite")
}

// TODO(trinhdn): continue testing other implemented methods
