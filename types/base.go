package types

import (
	"time"

	"github.com/kardiachain/go-kardia/types"
)

const (
	BloomByteLength = 256
)

type Bloom [BloomByteLength]byte

type Block struct {
	Hash   string `json:"hash,omitempty" bson:"hash"`
	Height uint64 `json:"height,omitempty" bson:"height"`

	CommitHash      string    `json:"commitHash,omitempty" bson:"commitHash"`
	GasLimit        uint64    `json:"gasLimit,omitempty" bson:"gasLimit"`
	GasUsed         uint64    `json:"gasUsed" bson:"gasUsed"`
	Rewards         string    `json:"rewards" bson:"rewards"`
	NumTxs          uint64    `json:"numTxs" bson:"numTxs"`
	Time            time.Time `json:"time,omitempty" bson:"time"`
	ProposerAddress string    `json:"proposerAddress,omitempty" bson:"proposerAddress"`

	LastBlock string `json:"lastBlock,omitempty" bson:"lastBlock"`

	DataHash     string `json:"dataHash,omitempty" bson:"dataHash"`
	ReceiptsRoot string `json:"receiptsRoot,omitempty" bson:"receiptsRoot"`
	LogsBloom    Bloom  `json:"logsBloom,omitempty" bson:"logsBloom"`

	ValidatorHash     string `json:"validatorHash,omitempty" bson:"validatorHash"`
	NextValidatorHash string `json:"nextValidatorHash,omitempty" bson:"nextValidatorHash"` // validators for the next block
	ConsensusHash     string `json:"consensusHash,omitempty" bson:"consensusHash"`
	AppHash           string `json:"appHash,omitempty" bson:"appHash"`
	EvidenceHash      string `json:"evidenceHash,omitempty" bson:"evidenceHash"`

	Txs      []*Transaction `json:"txs,omitempty" bson:"-"`
	Receipts []*Receipt     `json:"receipts,omitempty" bson:"-"`
}

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
	LogsBloom        Bloom         `json:"logsBloom"`
	Root             string        `json:"root"`
}

type Receipt struct {
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	TransactionHash   string `json:"transactionHash"`
	GasUsed           string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	ContractAddress   string `json:"contractAddress"`
	Logs              []*Log `json:"logs"`
	LogsBloom         string `json:"logsBloom"`
	Root              string `json:"root"`
	Status            string `json:"status"`
}

type Header struct {
	Hash              string         `json:"hash"`
	Height            uint64         `json:"height"`
	LastBlock         string         `json:"lastBlock"`
	LastBlockID       *types.BlockID `json:"lastBlockID"`
	CommitHash        string         `json:"commitHash"`
	Time              time.Time      `json:"time"`
	NumTxs            uint64         `json:"numTxs"`
	GasUsed           uint64         `json:"gasUsed"`
	GasLimit          uint64         `json:"gasLimit"`
	Rewards           string         `json:"Rewards"`
	ProposerAddress   string         `json:"proposerAddress"`
	TxHash            string         `json:"dataHash"`     // transactions
	ReceiptHash       string         `json:"receiptsRoot"` // receipt root
	Bloom             *types.Bloom   `json:"logsBloom"`
	ValidatorsHash    string         `json:"validatorHash"`     // current block validators hash
	NextValidatorHash string         `json:"nextValidatorHash"` // next block validators hash
	ConsensusHash     string         `json:"consensusHash"`     // current consensus hash
	AppHash           string         `json:"appHash"`           // state of transactions
	EvidenceHash      string         `json:"evidenceHash"`      // hash of evidence
}

type Log struct {
	Address       string                 `json:"address,omitempty" bson:"address"`
	MethodName    string                 `json:"methodName,omitempty" bson:"methodName"`
	ArgumentsName string                 `json:"argumentsName,omitempty" bson:"argumentsName"`
	Arguments     map[string]interface{} `json:"arguments,omitempty" bson:"arguments"`
	Topics        []string               `json:"topics,omitempty" bson:"topics"`
	Data          string                 `json:"data,omitempty" bson:"data"`
	BlockHeight   uint64                 `json:"blockHeight,omitempty" bson:"blockHeight"`
	Time          time.Time              `json:"time" bson:"time"`
	TxHash        string                 `json:"transactionHash"  bson:"transactionHash"`
	TxIndex       string                 `json:"transactionIndex,omitempty" bson:"transactionIndex"`
	BlockHash     string                 `json:"blockHash,omitempty" bson:"blockHash"`
	Index         string                 `json:"logIndex,omitempty" bson:"logIndex"`
	Removed       bool                   `json:"removed,omitempty" bson:"removed"`
}

type FunctionCall struct {
	Function   string                 `json:"function"`
	MethodID   string                 `json:"methodID"`
	MethodName string                 `json:"methodName"`
	Arguments  map[string]interface{} `json:"arguments"`
}
