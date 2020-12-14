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
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"go.uber.org/zap"

	"github.com/kardiachain/go-kardia"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/rpc"

	"github.com/kardiachain/go-kaiclient/types"
)

var (
	ErrParsingBigIntFromString = errors.New("cannot parse string to big.Int")
	ErrValidatorNotFound       = errors.New("validator address not found")

	tenPoweredBy5  = new(big.Int).Exp(big.NewInt(10), big.NewInt(5), nil)
	tenPoweredBy18 = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
)

type RPCClient struct {
	c      *rpc.Client
	isDead bool
	ip     string
}

// client return an *rpc.client instance
type client struct {
	clientList        []*RPCClient
	trustedClientList []*RPCClient
	defaultClient     *RPCClient
	numRequest        int
	lgr               *zap.Logger
}

// NewKaiClient creates a client that uses the given RPC client.
func NewKaiClient(cfg *Config) (Client, error) {
	if len(cfg.rpcURL) == 0 && len(cfg.trustedNodeRPCURL) == 0 {
		return nil, errors.New("empty RPC URL")
	}

	var (
		defaultClient *RPCClient = nil
		clientList               = []*RPCClient{}
	)
	for _, u := range cfg.rpcURL {
		rpcClient, err := rpc.Dial(u)
		if err != nil {
			return nil, err
		}
		newClient := &RPCClient{
			c:      rpcClient,
			isDead: false,
			ip:     u,
		}
		clientList = append(clientList, newClient)
	}
	var trustedClientList = []*RPCClient{}
	for _, u := range cfg.trustedNodeRPCURL {
		rpcClient, err := rpc.Dial(u)
		if err != nil {
			return nil, err
		}
		newClient := &RPCClient{
			c:      rpcClient,
			isDead: false,
			ip:     u,
		}
		trustedClientList = append(trustedClientList, newClient)
	}
	// set default RPC client as one of our trusted ones
	defaultClient = trustedClientList[0]

	return &client{clientList, trustedClientList, defaultClient, 0, cfg.lgr}, nil
}

func (ec *client) chooseClient() *RPCClient {
	if len(ec.clientList) > 1 {
		if ec.numRequest == len(ec.clientList)-1 {
			ec.numRequest = 0
		} else {
			ec.numRequest++
		}
		return ec.clientList[ec.numRequest%(len(ec.clientList)-1)]
	}
	return ec.defaultClient
}

// LatestBlockNumber gets latest block number
func (ec *client) LatestBlockNumber(ctx context.Context) (uint64, error) {
	var result uint64
	err := ec.defaultClient.c.CallContext(ctx, &result, "kai_blockNumber")
	return result, err
}

// BlockByHash returns the given full block.
//
// Use HeaderByHash if you don't need all transactions or uncle headers.
func (ec *client) BlockByHash(ctx context.Context, hash string) (*types.Block, error) {
	return ec.getBlock(ctx, "kai_getBlockByHash", common.HexToHash(hash))
}

// BlockByHeight returns a block from the current canonical chain.
//
// Use HeaderByNumber if you don't need all transactions or uncle headers.
// TODO(trinhdn): If number is nil, the latest known block is returned.
func (ec *client) BlockByHeight(ctx context.Context, height uint64) (*types.Block, error) {
	return ec.getBlock(ctx, "kai_getBlockByNumber", height)
}

// BlockHeaderByNumber returns a block header from the current canonical chain.
// TODO(trinhdn): If number is nil, the latest known block header is returned.
func (ec *client) BlockHeaderByNumber(ctx context.Context, number uint64) (*types.Header, error) {
	return ec.getBlockHeader(ctx, "kai_getBlockHeaderByNumber", number)
}

// BlockHeaderByHash returns the given block header.
func (ec *client) BlockHeaderByHash(ctx context.Context, hash string) (*types.Header, error) {
	return ec.getBlockHeader(ctx, "kai_getBlockHeaderByHash", common.HexToHash(hash))
}

// GetTransaction returns the transaction with the given hash.
func (ec *client) GetTransaction(ctx context.Context, hash string) (*types.Transaction, error) {
	var raw *types.Transaction
	err := ec.chooseClient().c.CallContext(ctx, &raw, "tx_getTransaction", common.HexToHash(hash))
	if err != nil {
		return nil, err
	} else if raw == nil {
		return nil, kardia.NotFound
	}
	return raw, nil
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (ec *client) GetTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error) {
	var r *types.Receipt
	err := ec.chooseClient().c.CallContext(ctx, &r, "tx_getTransactionReceipt", common.HexToHash(txHash))
	if err == nil {
		if r == nil {
			return nil, kardia.NotFound
		}
	}
	return r, err
}

// BalanceAt returns the wei balance of the given account.
// The block number can be nil, in which case the balance is taken from the latest known block.
func (ec *client) GetBalance(ctx context.Context, account string) (string, error) {
	var (
		result string
		err    error
	)
	err = ec.chooseClient().c.CallContext(ctx, &result, "account_balance", common.HexToAddress(account), "latest")
	return result, err
}

// StorageAt returns the value of key in the contract storage of the given account.
// The block number can be nil, in which case the value is taken from the latest known block.
func (ec *client) GetStorageAt(ctx context.Context, account string, key string) (common.Bytes, error) {
	var result common.Bytes
	err := ec.chooseClient().c.CallContext(ctx, &result, "kai_getStorageAt", common.HexToAddress(account), key, "latest")
	return result, err
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
func (ec *client) GetCode(ctx context.Context, account string) (common.Bytes, error) {
	var result common.Bytes
	err := ec.chooseClient().c.CallContext(ctx, &result, "kai_getCode", common.HexToAddress(account), "latest")
	return result, err
}

// NonceAt returns the account nonce of the given account.
func (ec *client) NonceAt(ctx context.Context, account string) (uint64, error) {
	var result uint64
	err := ec.chooseClient().c.CallContext(ctx, &result, "account_nonce", common.HexToAddress(account))
	return result, err
}

// SendRawTransaction injects a signed transaction into the pending pool for execution.
//
// If the transaction was a contract creation use the GetTransactionReceipt method to get the
// contract address after the transaction has been mined.
func (ec *client) SendRawTransaction(ctx context.Context, tx string) error {
	return ec.chooseClient().c.CallContext(ctx, nil, "tx_sendRawTransaction", tx)
}

func (ec *client) Peers(ctx context.Context, client *RPCClient) ([]*types.PeerInfo, error) {
	var result []*types.PeerInfo
	err := client.c.CallContext(ctx, &result, "node_peers")
	return result, err
}

func (ec *client) NodesInfo(ctx context.Context) ([]*types.NodeInfo, error) {
	var (
		nodes = []*types.NodeInfo(nil)
		err   error
	)
	clientList := append(ec.clientList, ec.trustedClientList...)
	nodeMap := make(map[string]*types.NodeInfo, len(clientList))
	for _, client := range clientList {
		var (
			node  *types.NodeInfo
			peers []*types.PeerInfo
		)
		err = client.c.CallContext(ctx, &node, "node_nodeInfo")
		if err != nil {
			continue
		}
		peers, err = ec.Peers(ctx, client)
		if err != nil {
			continue
		}
		node.Peers = peers
		nodeMap[node.ID] = node
	}
	for _, node := range nodeMap {
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (ec *client) Datadir(ctx context.Context) (string, error) {
	var result string
	err := ec.chooseClient().c.CallContext(ctx, &result, "node_datadir")
	return result, err
}

func (ec *client) Validator(ctx context.Context, address string) (*types.Validator, error) {
	var validator *types.Validator
	err := ec.defaultClient.c.CallContext(ctx, &validator, "kai_validator", address, true)
	if err != nil {
		return nil, err
	}
	return validator, nil
}

func (ec *client) Validators(ctx context.Context) (*types.Validators, error) {
	var validators []*types.Validator
	err := ec.defaultClient.c.CallContext(ctx, &validators, "kai_validators", true)
	if err != nil {
		return nil, err
	}
	var (
		delegators                 = make(map[string]bool)
		totalStakedAmount          = big.NewInt(0)
		totalDelegatorStakedAmount = big.NewInt(0)

		valStakedAmount *big.Int
		delStakedAmount *big.Int
		ok              bool
	)
	for _, val := range validators {
		for _, del := range val.Delegators {
			delegators[del.Address.Hex()] = true
			// exclude validator self delegation
			if del.Address.Equal(val.Address) {
				continue
			}
			delStakedAmount, ok = new(big.Int).SetString(del.StakedAmount, 10)
			if !ok {
				return nil, err
			}
			totalDelegatorStakedAmount = new(big.Int).Add(totalDelegatorStakedAmount, delStakedAmount)
		}
		valStakedAmount, ok = new(big.Int).SetString(val.StakedAmount, 10)
		if !ok {
			return nil, err
		}
		totalStakedAmount = new(big.Int).Add(totalStakedAmount, valStakedAmount)
	}
	sort.Slice(validators, func(i, j int) bool {
		iAmount, _ := new(big.Int).SetString(validators[i].StakedAmount, 10)
		jAmount, _ := new(big.Int).SetString(validators[j].StakedAmount, 10)
		return iAmount.Cmp(jAmount) == 1
	})
	for _, val := range validators {
		if val, err = convertValidatorInfo(val, totalStakedAmount); err != nil {
			return nil, err
		}
	}
	result := &types.Validators{
		TotalValidators:            len(validators),
		TotalDelegators:            len(delegators),
		TotalStakedAmount:          totalStakedAmount.String(),
		TotalValidatorStakedAmount: new(big.Int).Sub(totalStakedAmount, totalDelegatorStakedAmount).String(),
		TotalDelegatorStakedAmount: totalDelegatorStakedAmount.String(),
		TotalProposer:              21, // TODO(trinhdn): follow core API updates
		Validators:                 validators,
	}
	return result, nil
}

func (ec *client) getBlock(ctx context.Context, method string, args ...interface{}) (*types.Block, error) {
	var raw types.Block
	err := ec.defaultClient.c.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, err
	}
	return &raw, nil
}

func (ec *client) getBlockHeader(ctx context.Context, method string, args ...interface{}) (*types.Header, error) {
	var raw types.Header
	err := ec.defaultClient.c.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, err
	}
	return &raw, nil
}

func convertValidatorInfo(val *types.Validator, totalStakedAmount *big.Int) (*types.Validator, error) {
	var err error
	val.Commission = ""
	if val.CommissionRate, err = convertBigIntToPercentage(val.CommissionRate); err != nil {
		return nil, err
	}
	if val.MaxRate, err = convertBigIntToPercentage(val.MaxRate); err != nil {
		return nil, err
	}
	if val.MaxChangeRate, err = convertBigIntToPercentage(val.MaxChangeRate); err != nil {
		return nil, err
	}
	if totalStakedAmount != nil {
		if val.VotingPowerPercentage, err = calculateVotingPower(val.StakedAmount, totalStakedAmount); err != nil {
			return nil, err
		}
	}
	return val, nil
}

func convertBigIntToPercentage(raw string) (string, error) {
	input, ok := new(big.Int).SetString(raw, 10)
	if !ok {
		return "", ErrParsingBigIntFromString
	}
	tmp := new(big.Int).Mul(input, tenPoweredBy18)
	result := new(big.Int).Div(tmp, tenPoweredBy18).String()
	result = fmt.Sprintf("%020s", result)
	result = strings.TrimLeft(strings.TrimRight(strings.TrimRight(result[:len(result)-16]+"."+result[len(result)-16:], "0"), "."), "0")
	if strings.HasPrefix(result, ".") {
		result = "0" + result
	}
	return result, nil
}

func calculateVotingPower(raw string, total *big.Int) (string, error) {
	valStakedAmount, ok := new(big.Int).SetString(raw, 10)
	if !ok {
		return "", ErrParsingBigIntFromString
	}
	tmp := new(big.Int).Mul(valStakedAmount, tenPoweredBy5)
	result := new(big.Int).Div(tmp, total).String()
	result = fmt.Sprintf("%020s", result)
	result = strings.TrimLeft(strings.TrimRight(strings.TrimRight(result[:len(result)-3]+"."+result[len(result)-3:], "0"), "."), "0")
	if strings.HasPrefix(result, ".") {
		result = "0" + result
	}
	return result, nil
}
