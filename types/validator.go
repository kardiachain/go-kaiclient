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
// Package types
package types

import (
	"math/big"
	"time"

	"github.com/kardiachain/go-kardia/lib/common"
)

type Validators struct {
	TotalValidators            int          `json:"totalValidators"`
	TotalProposers             int          `json:"totalProposers"`
	TotalCandidates            int          `json:"totalCandidates"`
	TotalDelegators            int          `json:"totalDelegators"`
	TotalStakedAmount          string       `json:"totalStakedAmount"`
	TotalValidatorStakedAmount string       `json:"totalValidatorStakedAmount"`
	TotalDelegatorStakedAmount string       `json:"totalDelegatorStakedAmount"`
	Validators                 []*Validator `json:"validators"`
}

type Validator struct {
	Address               common.Address `json:"address"`
	SmcAddress            common.Address `json:"smcAddress"`
	Status                uint8          `json:"status"`
	Jailed                bool           `json:"jailed"`
	Name                  string         `json:"name,omitempty"`
	VotingPowerPercentage string         `json:"votingPowerPercentage"`
	StakedAmount          string         `json:"stakedAmount"`
	AccumulatedCommission string         `json:"accumulatedCommission"`
	UpdateTime            time.Time      `json:"updateTime"`
	CommissionRate        string         `json:"commissionRate"`
	TotalDelegators       int            `json:"totalDelegators"`
	MaxRate               string         `json:"maxRate"`
	MaxChangeRate         string         `json:"maxChangeRate"`
	MissedBlocks          uint64         `json:"missedBlocks"`
	IndicatorRate         float64        `json:"indicatorRate"`
	Delegators            []*Delegator   `json:"delegators,omitempty"`
}

type RPCValidator struct {
	Name                  [32]uint8       `json:"name"`
	ValAddr               common.Address  `json:"validatorAddress"`
	ValStakingSmc         common.Address  `json:"valStakingSmc"`
	Tokens                *big.Int        `json:"tokens"`
	Jailed                bool            `json:"jailed"`
	DelegationShares      *big.Int        `json:"delegationShares"`
	AccumulatedCommission *big.Int        `json:"accumulatedCommission"`
	UbdEntryCount         *big.Int        `json:"ubdEntryCount"`
	UpdateTime            *big.Int        `json:"updateTime"`
	Status                uint8           `json:"status"`
	UnbondingTime         *big.Int        `json:"unbondingTime"`
	UnbondingHeight       *big.Int        `json:"unbondingHeight"`
	CommissionRate        *big.Int        `json:"commissionRate,omitempty"`
	MaxRate               *big.Int        `json:"maxRate,omitempty"`
	MaxChangeRate         *big.Int        `json:"maxChangeRate,omitempty"`
	Delegators            []*RPCDelegator `json:"delegators,omitempty"`
}

type RPCDelegator struct {
	Address      common.Address `json:"address"`
	StakedAmount *big.Int       `json:"stakedAmount"`
	Reward       *big.Int       `json:"reward"`
}

type Delegator struct {
	Address      common.Address `json:"address"`
	Name         string         `json:"name,omitempty"`
	StakedAmount string         `json:"stakedAmount"`
	Reward       string         `json:"reward"`
}

type SlashEvents struct {
	Period   string `json:"period"`
	Fraction string `json:"fraction"`
	Height   string `json:"height"`
}

type SigningInfo struct {
	StartHeight        uint64
	IndexOffset        uint64
	Tombstoned         bool
	MissedBlockCounter uint64
	JailedUntil        uint64
}

type RPCPeerInfo struct {
	NodeInfo         *NodeInfo `json:"node_info"`
	IsOutbound       bool      `json:"is_outbound"`
	ConnectionStatus struct {
		Duration uint64 `json:"Duration"`
	} `json:"connection_status"`
	RemoteIP string `json:"remote_ip"`
}

type PeerInfo struct {
	Duration uint64 `json:"Duration,omitempty"`
	Moniker  string `json:"moniker,omitempty"` // arbitrary moniker
	RemoteIP string `json:"remote_ip,omitempty"`
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
	Peers           map[string]*PeerInfo `json:"peers,omitempty"` // peers details
	Other           DefaultNodeInfoOther `json:"other"`           // other application specific data
}

type ValidatorsByDelegator struct {
	Name                  string         `json:"name"`
	Validator             common.Address `json:"validator"`
	ValidatorContractAddr common.Address `json:"validatorContractAddr"`
	ValidatorStatus       uint8          `json:"validatorStatus"`
	StakedAmount          string         `json:"stakedAmount"`
	ClaimableRewards      string         `json:"claimableRewards"`
	UnbondedAmount        string         `json:"unbondedAmount"`
	WithdrawableAmount    string         `json:"withdrawableAmount"`
}
