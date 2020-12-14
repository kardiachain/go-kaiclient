package types

import (
	"time"

	"github.com/kardiachain/go-kardia/types"
)

type Transaction struct {
	BlockHash   string `json:"blockHash" bson:"blockHash"`
	BlockNumber uint64 `json:"blockNumber" bson:"blockNumber"`

	Hash             string      `json:"hash" bson:"hash"`
	From             string      `json:"from" bson:"from"`
	To               string      `json:"to" bson:"to"`
	Status           uint        `json:"status" bson:"status"`
	ContractAddress  string      `json:"contractAddress" bson:"contractAddress"`
	Value            string      `json:"value" bson:"value"`
	GasPrice         uint64      `json:"gasPrice" bson:"gasPrice"`
	GasLimit         uint64      `json:"gas" bson:"gas"`
	GasUsed          uint64      `json:"gasUsed"`
	TxFee            string      `json:"txFee"`
	Nonce            uint64      `json:"nonce" bson:"nonce"`
	Time             time.Time   `json:"time" bson:"time"`
	InputData        string      `json:"input" bson:"input"`
	Logs             []Log       `json:"logs" bson:"logs"`
	TransactionIndex uint        `json:"transactionIndex"`
	LogsBloom        types.Bloom `json:"logsBloom"`
	Root             string      `json:"root"`
}

type Log struct {
	Address     string   `json:"address"`
	Topics      []string `json:"topics"`
	Data        string   `json:"data"`
	BlockHeight uint64   `json:"blockHeight"`
	TxHash      string   `json:"transactionHash"`
	TxIndex     uint     `json:"transactionIndex"`
	BlockHash   string   `json:"blockHash"`
	Index       uint     `json:"logIndex"`
	Removed     bool     `json:"removed"`
}

type Receipt struct {
	BlockHash   string `json:"blockHash"`
	BlockHeight uint64 `json:"blockHeight"`

	TransactionHash  string `json:"transactionHash"`
	TransactionIndex uint64 `json:"transactionIndex"`

	From              string      `json:"from"`
	To                string      `json:"to"`
	GasUsed           uint64      `json:"gasUsed"`
	CumulativeGasUsed uint64      `json:"cumulativeGasUsed"`
	ContractAddress   string      `json:"contractAddress"`
	Logs              []Log       `json:"logs"`
	LogsBloom         types.Bloom `json:"logsBloom"`
	Root              string      `json:"root"`
	Status            uint        `json:"status"`
}

type TransactionByAddress struct {
	Address string    `json:"address" bson:"address"`
	TxHash  string    `json:"txHash" bson:"txHash"`
	Time    time.Time `json:"time" bson:"time"`
}
