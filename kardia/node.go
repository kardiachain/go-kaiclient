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
	"os"
	"path"
	"runtime"

	"github.com/kardiachain/go-kardia"
	"github.com/kardiachain/go-kardia/lib/abi/bind"
	"go.uber.org/zap"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/rpc"
)

const (
	StakingContractAddr = "0x0000000000000000000000000000000000001337"
)

type Node interface {
	IsAlive() bool
	Info(ctx context.Context) (*NodeInfo, error)

	IAddress

	LatestBlockNumber(ctx context.Context) (uint64, error)
	BlockByHash(ctx context.Context, hash string) (*Block, error)
	BlockByHeight(ctx context.Context, height uint64) (*Block, error)
	BlockHeaderByHash(ctx context.Context, hash string) (*Header, error)
	BlockHeaderByNumber(ctx context.Context, number uint64) (*Header, error)

	DecodeInputData(to string, input string) (*FunctionCall, error)

	IContract
	IStaking
	ITx

	GetCirculatingSupply(ctx context.Context) (*big.Int, error)

	KardiaCall(ctx context.Context, args SMCCallArgs) ([]byte, error)
	IValidator
	IDelegator
	bind.ContractCaller
	bind.ContractTransactor
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
}

func (n *node) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	panic("implement me")
}

func (n *node) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	panic("implement me")
}

func (n *node) SuggestGasPrice(ctx context.Context) (uint64, error) {
	panic("implement me")
}

func (n *node) EstimateGas(ctx context.Context, call kardia.CallMsg) (gas uint64, err error) {
	panic("implement me")
}

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
		lgr:    lgr,
	}
	if err := node.setupSMC(); err != nil {
		return nil, err
	}
	return node, nil
}

func (n *node) setupSMC() error {
	filePath := os.Getenv("ABI_PATH")
	if filePath == "" {
		_, filename, _, _ := runtime.Caller(1)
		filePath = path.Dir(filename)
	}

	stakingABI, err := os.Open(path.Join(filePath, "../kardia/abi/staking.json"))
	if err != nil {
		panic("cannot read staking ABI file")
	}
	stakingSmcABI, err := abi.JSON(stakingABI)
	if err != nil {
		return err
	}
	stakingUtil := &Contract{
		Abi:             &stakingSmcABI,
		ContractAddress: common.HexToAddress(StakingContractAddr),
	}
	n.stakingSMC = stakingUtil
	validatorABI, err := os.Open(path.Join(filePath, "../kardia/abi/validator.json"))
	if err != nil {
		return err
	}
	validatorSmcAbi, err := abi.JSON(validatorABI)
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
	paramsABI, err := os.Open(path.Join(filePath, "../kardia/abi/params.json"))
	if err != nil {
		return err
	}
	paramsSmcAbi, err := abi.JSON(paramsABI)
	if err != nil {
		return err
	}
	paramsUtil := &Contract{
		Abi:             &paramsSmcAbi,
		ContractAddress: paramsSmcAddr,
	}
	n.paramsSMC = paramsUtil

	return nil
}

func (n *node) IsAlive() bool {
	return true
}

func (n *node) Info(ctx context.Context) (*NodeInfo, error) {
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
