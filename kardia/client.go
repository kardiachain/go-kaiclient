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
// Package kardia
package kardia

import (
	"context"
	"math/big"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/mainchain/staking"

	"github.com/kardiachain/go-kaiclient/types"
)

type ClientInterface interface {
	LatestBlockNumber(ctx context.Context) (uint64, error)
	BlockByHash(ctx context.Context, hash string) (*types.Block, error)
	BlockByHeight(ctx context.Context, height uint64) (*types.Block, error)
	BlockHeaderByHash(ctx context.Context, hash string) (*types.Header, error)
	BlockHeaderByNumber(ctx context.Context, number uint64) (*types.Header, error)
	GetTransaction(ctx context.Context, hash string) (*types.Transaction, error)
	GetTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error)
	GetBalance(ctx context.Context, account string) (string, error)
	GetStorageAt(ctx context.Context, account string, key string) (common.Bytes, error)
	GetCode(ctx context.Context, account string) (common.Bytes, error)
	NonceAt(ctx context.Context, account string) (uint64, error)
	SendRawTransaction(ctx context.Context, tx string) error
	KardiaCall(ctx context.Context, args types.CallArgsJSON) (common.Bytes, error)
	Peers(ctx context.Context, client *RPCClient) ([]*types.PeerInfo, error)
	NodesInfo(ctx context.Context) ([]*types.NodeInfo, error)
	Datadir(ctx context.Context) (string, error)
	Validator(ctx context.Context, address string) (*types.Validator, error)
	Validators(ctx context.Context) (*types.Validators, error)

	// staking related methods
	GetValidatorsByDelegator(ctx context.Context, delAddr common.Address) ([]*types.ValidatorsByDelegator, error)

	// validator related methods
	GetValidatorParams(ctx context.Context, valSmcAddr common.Address) (*types.ValidatorParams, error)
	GetValidatorInfo(ctx context.Context, valSmcAddr common.Address) (*staking.Validator, error)
	GetDelegationRewards(ctx context.Context, valSmcAddr common.Address, delegatorAddr common.Address) (*big.Int, error)
	GetDelegatorStakedAmount(ctx context.Context, valSmcAddr common.Address, delegatorAddr common.Address) (*big.Int, error)
	GetUDBEntries(ctx context.Context, valSmcAddr common.Address, delegatorAddr common.Address) (*big.Int, *big.Int, error)
}

type Config struct {
	rpcURL        []string
	defaultRPCURL []string
	contracts     map[string]SmcUtil
}

type SmcUtil struct {
	Abi             *abi.ABI
	ContractAddress common.Address
	Bytecode        string
}

func NewConfig(rpcURL []string, defaultRPCURL []string, contracts map[string]SmcUtil) *Config {
	return &Config{
		rpcURL:        rpcURL,
		defaultRPCURL: defaultRPCURL,
		contracts:     contracts,
	}
}
