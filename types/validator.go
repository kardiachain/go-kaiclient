// Package types
package types

import (
	"github.com/kardiachain/go-kardiamain/lib/common"
)

type Validators struct {
	TotalValidators            int          `json:"totalValidators"`
	TotalDelegators            int          `json:"totalDelegators"`
	TotalStakedAmount          string       `json:"totalStakedAmount"`
	TotalValidatorStakedAmount string       `json:"totalValidatorStakedAmount"`
	TotalDelegatorStakedAmount string       `json:"totalDelegatorStakedAmount"`
	TotalProposer              int          `json:"totalProposer"`
	Validators                 []*Validator `json:"validators"`
}

type Validator struct {
	Address               common.Address `json:"address"`
	Name                  string         `json:"name,omitempty"`
	VotingPower           int64          `json:"votingPower"`
	VotingPowerPercentage string         `json:"votingPowerPercentage"`
	StakedAmount          string         `json:"stakedAmount"`
	Commission            string         `json:"commission"`
	CommissionRate        string         `json:"commissionRate"`
	TotalDelegators       int            `json:"totalDelegators"`
	MaxRate               string         `json:"maxRate"`
	MaxChangeRate         string         `json:"maxChangeRate"`
	Delegators            []*Delegator   `json:"delegators,omitempty"`
}

type Delegator struct {
	Address      common.Address `json:"address"`
	Name         string         `json:"name,omitempty"`
	StakedAmount string         `json:"stakedAmount"`
	Reward       string         `json:"reward"`
}

type PeerInfo struct {
	NodeInfo         *NodeInfo `json:"node_info"`
	IsOutbound       bool      `json:"is_outbound"`
	ConnectionStatus struct {
		Duration uint64 `json:"Duration"`
	} `json:"connection_status"`
	RemoteIP string `json:"remote_ip"`
}

type ProtocolVersion struct {
	P2P   uint64 `json:"p2p"`
	Block uint64 `json:"block"`
	App   uint64 `json:"app"`
}

type DefaultNodeInfoOther struct {
	TxIndex    string `json:"tx_index"`
	RPCAddress string `json:"rpc_address"`
}

type NodeInfo struct {
	ProtocolVersion ProtocolVersion      `json:"protocol_version"`
	ID              string               `json:"id"`              // authenticated identifier
	ListenAddr      string               `json:"listen_addr"`     // accepting incoming
	Network         string               `json:"network"`         // network/chain ID
	Version         string               `json:"version"`         // major.minor.revision
	Moniker         string               `json:"moniker"`         // arbitrary moniker
	Peers           []*PeerInfo          `json:"peers,omitempty"` // peers details
	Other           DefaultNodeInfoOther `json:"other"`           // other application specific data
}
