/*
 *  Copyright 2020 KardiaChain
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
package types

import (
	"math/big"
	"time"

	"github.com/kardiachain/go-kardia/types"
)

type Transaction struct {
	BlockHash   string `json:"blockHash" bson:"blockHash"`
	BlockNumber uint64 `json:"blockNumber" bson:"blockNumber"`

	Hash             string        `json:"hash" bson:"hash"`
	From             string        `json:"from" bson:"from"`
	To               string        `json:"to" bson:"to"`
	Status           uint          `json:"status" bson:"status"`
	ContractAddress  string        `json:"contractAddress" bson:"contractAddress"`
	Value            string        `json:"value" bson:"value"`
	GasPrice         uint64        `json:"gasPrice" bson:"gasPrice"`
	GasLimit         uint64        `json:"gas" bson:"gas"`
	GasUsed          uint64        `json:"gasUsed"`
	TxFee            string        `json:"txFee"`
	Nonce            uint64        `json:"nonce" bson:"nonce"`
	Time             time.Time     `json:"time" bson:"time"`
	InputData        string        `json:"input" bson:"input"`
	DecodedInputData *FunctionCall `json:"decodedInputData,omitempty" bson:"decodedInputData"`
	Logs             []Log         `json:"logs" bson:"logs"`
	TransactionIndex uint          `json:"transactionIndex"`
	LogsBloom        types.Bloom   `json:"logsBloom"`
	Root             string        `json:"root"`
}

type FunctionCall struct {
	Function   string                 `json:"function"`
	MethodID   string                 `json:"methodID"`
	MethodName string                 `json:"methodName"`
	Arguments  map[string]interface{} `json:"arguments"`
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

type CallArgsJSON struct {
	From     string   `json:"from"`     // the sender of the 'transaction'
	To       *string  `json:"to"`       // the destination contract (nil for contract creation)
	Gas      uint64   `json:"gas"`      // if 0, the call executes with near-infinite gas
	GasPrice *big.Int `json:"gasPrice"` // wei <-> gas exchange ratio
	Value    *big.Int `json:"value"`    // amount of wei sent along with the call
	Data     string   `json:"data"`     // input data, usually an ABI-encoded contract method invocation
}
