package types

import (
	"github.com/kardiachain/go-kardia/lib/common"
)

type Filter struct {
	Address   common.Address `json:"address,omitempty"`
	FromBlock string         `json:"fromBlock,omitempty"`
	ToBlock   string         `json:"toBlock,omitempty"`
	Topics    []string       `json:"topics,omitempty"`
}

type Event struct {
	Address         common.Address `json:"address"`
	BlockHash       common.Hash    `json:"blockHash"`
	BlockNumber     string         `json:"blockNumber"`
	Topics          []string       `json:"topics,omitempty"`
	TransactionHash common.Hash    `json:"transactionHash"`
	Data            string         `json:"data,omitempty"`
}

type EventData struct {
	Address         common.Address `json:"address"`
	BlockHash       common.Hash    `json:"blockHash"`
	BlockNumber     string         `json:"blockNumber"`
	Topics          [][]byte       `json:"topics,omitempty"`
	TransactionHash common.Hash    `json:"transactionHash"`
	Data            []byte         `json:"data,omitempty"`
}
