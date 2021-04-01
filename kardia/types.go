// Package kardia
package kardia

import (
	"math/big"
	"time"

	"github.com/kardiachain/go-kardia/lib/common"
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
	TransactionHash   string      `json:"transactionHash"`
	GasUsed           uint64      `json:"gasUsed"`
	CumulativeGasUsed uint64      `json:"cumulativeGasUsed"`
	ContractAddress   string      `json:"contractAddress"`
	Logs              []*Log      `json:"logs"`
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
	Address       string                 `json:"address,omitempty" bson:"address"`
	MethodName    string                 `json:"methodName,omitempty" bson:"methodName"`
	ArgumentsName string                 `json:"argumentsName,omitempty" bson:"argumentsName"`
	Arguments     map[string]interface{} `json:"arguments,omitempty" bson:"arguments"`
	Topics        []string               `json:"topics,omitempty" bson:"topics"`
	Data          string                 `json:"data,omitempty" bson:"data"`
	BlockHeight   uint64                 `json:"blockHeight,omitempty" bson:"blockHeight"`
	Time          time.Time              `json:"time" bson:"time"`
	TxHash        string                 `json:"transactionHash"  bson:"transactionHash"`
	TxIndex       uint                   `json:"transactionIndex,omitempty" bson:"transactionIndex"`
	BlockHash     string                 `json:"blockHash,omitempty" bson:"blockHash"`
	Index         uint                   `json:"logIndex,omitempty" bson:"logIndex"`
	Removed       bool                   `json:"removed,omitempty" bson:"removed"`
}

type Header struct {
	Hash              string    `json:"hash"`
	Height            uint64    `json:"height"`
	LastBlock         string    `json:"lastBlock"`
	CommitHash        string    `json:"commitHash"`
	Time              time.Time `json:"time"`
	NumTxs            uint64    `json:"numTxs"`
	GasUsed           uint64    `json:"gasUsed"`
	GasLimit          uint64    `json:"gasLimit"`
	Rewards           string    `json:"Rewards"`
	ProposerAddress   string    `json:"proposerAddress"`
	TxHash            string    `json:"dataHash"`     // transactions
	ReceiptHash       string    `json:"receiptsRoot"` // receipt root
	Bloom             string    `json:"logsBloom"`
	ValidatorsHash    string    `json:"validatorHash"`     // current block validators hash
	NextValidatorHash string    `json:"nextValidatorHash"` // next block validators hash
	ConsensusHash     string    `json:"consensusHash"`     // current consensus hash
	AppHash           string    `json:"appHash"`           // state of transactions
	EvidenceHash      string    `json:"evidenceHash"`      // hash of evidence
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

type Event struct {
	Name string
	// RawName is the raw event name parsed from ABI.
	RawName    string
	Inputs     map[string]interface{}
	TxHash     string
	SMCAddress string
}

type FilterArgs struct {
	From    uint64
	To      uint64
	Address []string
	Topics  []string
}

type FilterLogs struct {
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockHeight      uint64   `json:"blockHeight"`
	Data             string   `json:"data"`
	LogIndex         int64    `json:"logIndex"`
	Removed          bool     `json:"removed"`
	Topic            []string `json:"topic"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex int64    `json:"transactionIndex"`
}

type ValidatorsByDelegator struct {
	Name                    string            `json:"name"`
	Validator               common.Address    `json:"validator"`
	ValidatorContractAddr   common.Address    `json:"validatorContractAddr"`
	ValidatorStatus         uint8             `json:"validatorStatus"`
	ValidatorRole           int               `json:"validatorRole"`
	StakedAmount            string            `json:"stakedAmount"`
	ClaimableRewards        string            `json:"claimableRewards"`
	UnbondedRecords         []*UnbondedRecord `json:"unbondedRecords"`
	TotalWithdrawableAmount string            `json:"totalWithdrawableAmount"`
	TotalUnbondedAmount     string            `json:"totalUnbondedAmount"`
	UnbondedAmount          string            `json:"unbondedAmount"`
	WithdrawableAmount      string            `json:"withdrawableAmount"`
}

type UnbondedRecord struct {
	Balances        []*big.Int `json:"balances"`
	CompletionTimes []*big.Int `json:"completionTimes"`
}

type KRC20 struct {
	Address     common.Address
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply *big.Int
}
