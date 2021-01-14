// Package kardia
package kardia

import (
	"math/big"
	"time"

	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/types"
)

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
	LogsBloom    string `json:"logsBloom,omitempty" bson:"logsBloom"`

	ValidatorHash     string `json:"validatorHash,omitempty" bson:"validatorHash"`
	NextValidatorHash string `json:"nextValidatorHash,omitempty" bson:"nextValidatorHash"` // validators for the next block
	ConsensusHash     string `json:"consensusHash,omitempty" bson:"consensusHash"`
	AppHash           string `json:"appHash,omitempty" bson:"appHash"`
	EvidenceHash      string `json:"evidenceHash,omitempty" bson:"evidenceHash"`

	// Dual nodes
	NumDualEvents  uint64 `json:"numDualEvents,omitempty" bson:"numDualEvents"`
	DualEventsHash string `json:"dualEventsHash,omitempty" bson:"dualEventsHash"`

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
	LogsBloom        types.Bloom   `json:"logsBloom"`
	Root             string        `json:"root"`
}

type Validator struct {
	Name       [32]uint8
	Signer     common.Address
	SMCAddress common.Address
	Tokens     *big.Int
	Jailed     bool

	DelegationShares      *big.Int
	AccumulatedCommission *big.Int
	UbdEntryCount         *big.Int
	UpdateTime            *big.Int
	MinSelfDelegation     *big.Int
	Status                uint8

	UnbondingTime   *big.Int
	UnbondingHeight *big.Int

	Commission  *Commission
	SigningInfo *SigningInfo
	Delegators  []*Delegator
}

type Commission struct {
	Rate          *big.Int
	MaxRate       *big.Int
	MaxChangeRate *big.Int
}

type Delegator struct {
	Address      common.Address
	StakedAmount *big.Int
	Reward       *big.Int
}

type SigningInfo struct {
	StartHeight        *big.Int
	IndexOffset        *big.Int
	Tombstoned         bool
	MissedBlockCounter *big.Int
	JailedUntil        *big.Int
}

type PeerInfo struct {
	IsOutbound       bool `json:"is_outbound"`
	ConnectionStatus struct {
		Duration uint64 `json:"Duration"`
	} `json:"connection_status"`
	RemoteIP string `json:"remote_ip"`
}

type SMCCallArgs struct {
	From     string   `json:"from"`     // the sender of the 'transaction'
	To       *string  `json:"to"`       // the destination contract (nil for contract creation)
	Gas      uint64   `json:"gas"`      // if 0, the call executes with near-infinite gas
	GasPrice *big.Int `json:"gasPrice"` // HYDRO <-> gas exchange ratio
	Value    *big.Int `json:"value"`    // amount of HYDRO sent along with the call
	Data     string   `json:"data"`     // input data, usually an ABI-encoded contract method invocation
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

type FunctionCall struct {
	Function   string                 `json:"function"`
	MethodID   string                 `json:"methodID"`
	MethodName string                 `json:"methodName"`
	Arguments  map[string]interface{} `json:"arguments"`
}

type Log struct {
	Address     string `json:"address"`
	Name        string `json:"name"`
	MethodName  string `json:"methodName"`
	Params      []string
	Arguments   map[string]interface{} `json:"arguments"`
	Topics      []string               `json:"topics"`
	Data        string                 `json:"data"`
	BlockHeight uint64                 `json:"blockHeight"`
	TxHash      string                 `json:"transactionHash"`
	TxIndex     uint                   `json:"transactionIndex"`
	BlockHash   string                 `json:"blockHash"`
	Index       uint                   `json:"logIndex"`
	Removed     bool                   `json:"removed"`
}

type Header struct {
	Hash   string `json:"hash" bson:"blockHash"`
	Height uint64 `json:"height" bson:"height"`

	CommitHash      string    `json:"commitHash" bson:"commitHash"`
	GasLimit        uint64    `json:"gasLimit" bson:"gasLimit"`
	GasUsed         uint64    `json:"gasUsed" bson:"gasUsed"`
	Rewards         string    `json:"rewards" bson:"rewards"`
	NumTxs          uint64    `json:"numTxs" bson:"numTxs"`
	Time            time.Time `json:"time" bson:"time"`
	ProposerAddress string    `json:"proposerAddress" bson:"proposerAddress"`

	LastBlock string `json:"lastBlock" bson:"lastBlock"`

	DataHash     string `json:"dataHash" bson:"dataHash"`
	ReceiptsRoot string `json:"receiptsRoot" bson:"receiptsRoot"`
	LogsBloom    string `json:"logsBloom" bson:"logsBloom"`

	ValidatorHash     string `json:"validatorHash" bson:"validatorHash"`
	NextValidatorHash string `json:"nextValidatorHash" bson:"nextValidatorHash"` // validators for the next block
	ConsensusHash     string `json:"consensusHash" bson:"consensusHash"`
	AppHash           string `json:"appHash" bson:"appHash"`
	EvidenceHash      string `json:"evidenceHash" bson:"evidenceHash"`
}

type NodeInfo struct {
	ProtocolVersion ProtocolVersion `json:"protocol_version" bson:"protocolVersion"`
	ID              string          `json:"id" bson:"id"`                  // authenticated identifier
	ListenAddr      string          `json:"listen_addr" bson:"listenAddr"` // accepting incoming
	Network         string          `json:"network" bson:"network"`        // network/chain ID
	Version         string          `json:"version" bson:"version"`        // major.minor.revision
	Moniker         string          `json:"moniker" bson:"moniker"`        // arbitrary moniker
	Peers           []*PeerInfo     `json:"peers,omitempty" bson:"peers"`  // peers details
}

type ProtocolVersion struct {
	P2P   uint64 `json:"p2p"`
	Block uint64 `json:"block"`
	App   uint64 `json:"app"`
}

type SlashEvents struct {
	Period   string `json:"period" bson:"period"`
	Fraction string `json:"fraction" bson:"fraction"`
	Height   string `json:"height" bson:"height"`
}
