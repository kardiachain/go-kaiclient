package node

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kardiachain/go-kaiclient/rpc"
)

func TestGetBlockByNumber(t *testing.T) {
	c, err := rpc.NewClient("https://rpc.kardiachain.io")
	assert.Nil(t, err)

	eth := NewEth(c)
	blockNumber, err := eth.GetBlockNumber()
	assert.Nil(t, err)

	fmt.Printf("Block: %d \n", blockNumber)
	block, err := eth.GetBlockByNumber(big.NewInt(int64(blockNumber)), true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("block hash %s has %v txs\n", block.Hash, len(block.Txs))
}

func TestPollBlock(t *testing.T) {
	c, err := rpc.NewClientWithProxy("https://rpc.flashbots.net", "http://127.0.0.1:7890")
	if err != nil {
		t.Fatal(err)
	}
	eth := NewEth(c)
	for {
		blockNumber, err := eth.GetBlockNumber()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("get block %v\n", blockNumber)
		block, err := eth.GetBlockByNumber(big.NewInt(int64(blockNumber)), true)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("block hash %s has %v txs\n", block.Hash, len(block.Txs))
		time.Sleep(time.Duration(5) * time.Second)
	}

}
