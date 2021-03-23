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
	"math/big"
	"strings"

	"github.com/kardiachain/go-kardia"
	"github.com/kardiachain/go-kardia/lib/abi/bind"
	"github.com/kardiachain/go-kardia/lib/event"
	"github.com/kardiachain/go-kardia/types"
	"go.uber.org/zap"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/rpc"

	"github.com/kardiachain/go-kaiclient/kardia/smc"
)

const (
	StakingContractAddr = "0x0000000000000000000000000000000000001337"
)

type Node interface {
	Url() string
	IsAlive() bool
	NodeInfo(ctx context.Context) (*NodeInfo, error)
	GetCirculatingSupply(ctx context.Context) (*big.Int, error)
	KardiaCall(ctx context.Context, args SMCCallArgs) ([]byte, error)

	IAddress
	IBlock
	IReceipt
	IContract
	IStaking
	ITx
	ISubscription
	IToken

	IValidator
	IDelegator

	bind.ContractCaller
	bind.ContractTransactor
	bind.ContractBackend
}

func (n *node) FilterLogs(ctx context.Context, query kardia.FilterQuery) ([]types.Log, error) {
	panic("implement me")
}

func (n *node) SubscribeFilterLogs(ctx context.Context, query kardia.FilterQuery, ch chan<- types.Log) (event.Subscription, error) {
	panic("implement me")
}

type node struct {
	client *rpc.Client
	isLive bool
	url    string

	lgr *zap.Logger

	// SMC
	stakingSMC   *Contract
	validatorSMC *Contract
	paramsSMC    *Contract
	krc20SMC     *Contract
	krc721SMC    *Contract
}

func (n *node) Url() string {
	return n.url
}

//ContractTransactor
func (n *node) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	panic("implement me")
}

func (n *node) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return n.NonceAt(ctx, account.String())
}

func (n *node) SuggestGasPrice(ctx context.Context) (uint64, error) {
	panic("implement me")
}

func (n *node) EstimateGas(ctx context.Context, call kardia.CallMsg) (gas uint64, err error) {
	panic("implement me")
}

//ContractCaller
func (n *node) CodeAt(ctx context.Context, contract common.Address, blockNumber uint64) ([]byte, error) {
	panic("implement me")
}

func (n *node) CallContract(ctx context.Context, call kardia.CallMsg, blockNumber uint64) ([]byte, error) {
	panic("implement me")
}

func NewNode(url string, lgr *zap.Logger) (Node, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	node := &node{
		client: rpcClient,
		url:    url,
		lgr:    lgr,
	}
	if err := node.setupSMC(); err != nil {
		return nil, err
	}
	return node, nil
}

func (n *node) setupSMC() error {
	stakingSmcABI, err := abi.JSON(strings.NewReader(smc.StakingABI))
	if err != nil {
		return err
	}
	stakingUtil := &Contract{
		Abi:             &stakingSmcABI,
		ContractAddress: common.HexToAddress(StakingContractAddr),
	}
	n.stakingSMC = stakingUtil
	validatorSmcAbi, err := abi.JSON(strings.NewReader(smc.ValidatorABI))
	if err != nil {
		return err
	}
	validatorUtil := &Contract{
		Abi: &validatorSmcAbi,
	}
	n.validatorSMC = validatorUtil
	paramsSmcAddr, err := getParamsSMCAddress(stakingUtil, n.client)
	if err != nil {
		return err
	}
	paramsSmcAbi, err := abi.JSON(strings.NewReader(smc.ParamsABI))
	if err != nil {
		return err
	}
	paramsUtil := &Contract{
		Abi:             &paramsSmcAbi,
		ContractAddress: paramsSmcAddr,
	}
	n.paramsSMC = paramsUtil

	krc20SmcABI, err := abi.JSON(strings.NewReader(smc.KRC20ABI))
	if err != nil {
		return err
	}
	krc20Util := &Contract{
		Abi: &krc20SmcABI,
	}
	n.krc20SMC = krc20Util

	//stakingSmcABI, err := abi.JSON(strings.NewReader(smc.StakingABI))
	//if err != nil {
	//	return err
	//}
	//stakingUtil := &Contract{
	//	Abi:             &stakingSmcABI,
	//	ContractAddress: common.HexToAddress(StakingContractAddr),
	//}
	//n.stakingSMC = stakingUtil

	return nil
}

func (n *node) IsAlive() bool {
	return true
}

func (n *node) NodeInfo(ctx context.Context) (*NodeInfo, error) {
	var (
		node  *NodeInfo
		peers []*PeerInfo
	)
	// get current node info then get it's peers
	if err := n.client.CallContext(ctx, &node, "node_nodeInfo"); err != nil {
		return nil, err
	}

	err := n.client.CallContext(ctx, &peers, "node_peers")
	if err != nil {
		return nil, err
	}

	for _, p := range peers {
		node.Peers = append(node.Peers, &PeerInfo{
			//Moniker:  p.NodeInfo.Moniker,
			//Duration: p.ConnectionStatus.Duration,
			RemoteIP: p.RemoteIP,
		})
	}

	return node, nil
}

func (n *node) KardiaCall(ctx context.Context, args SMCCallArgs) ([]byte, error) {
	var result common.Bytes
	err := n.client.CallContext(ctx, &result, "kai_kardiaCall", args, "latest")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (n *node) DeployContract(ctx context.Context, auth *bind.TransactOpts, contract *BoundContract) {

}
