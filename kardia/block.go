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

	"github.com/kardiachain/go-kardia"
	"github.com/kardiachain/go-kardia/lib/common"
	kai "github.com/kardiachain/go-kardia/mainchain"
	"github.com/kardiachain/go-kardia/types"
)

type IBlock interface {
	LatestBlockNumber(ctx context.Context) (uint64, error)
	BlockByHash(ctx context.Context, hash string) (*Block, error)
	BlockByHeight(ctx context.Context, height uint64) (*Block, error)
	BlockHeaderByHash(ctx context.Context, hash string) (*Header, error)
	BlockHeaderByNumber(ctx context.Context, number uint64) (*Header, error)

	GetValidators(ctx context.Context, height uint64) (*types.ValidatorSet, error)
	GetCommit(ctx context.Context, height uint64) (*types.Commit, error)
	GetProof(ctx context.Context, address common.Address, storageKeys []string, height uint64) (*kai.AccountResult, error)
	FullHeaderByNumber(ctx context.Context, height uint64) (*FullHeader, error)
}

// BlockByHash returns the given full block.
// Use HeaderByHash if you don't need all transactions or uncle headers.
func (n *node) BlockByHash(ctx context.Context, hash string) (*Block, error) {
	return n.getBlock(ctx, "kai_getBlockByHash", common.HexToHash(hash))
}

// BlockByHeight returns a block from the current canonical chain.
// Use HeaderByNumber if you don't need all transactions or uncle headers.
func (n *node) BlockByHeight(ctx context.Context, height uint64) (*Block, error) {
	return n.getBlock(ctx, "kai_getBlockByNumber", height)
}

// BlockHeaderByNumber returns a block header from the current canonical chain.
func (n *node) BlockHeaderByNumber(ctx context.Context, number uint64) (*Header, error) {
	return n.getBlockHeader(ctx, "kai_getBlockHeaderByNumber", number)
}

// BlockHeaderByHash returns the given block header.
func (n *node) BlockHeaderByHash(ctx context.Context, hash string) (*Header, error) {
	return n.getBlockHeader(ctx, "kai_getBlockHeaderByHash", common.HexToHash(hash))
}

// LatestBlockNumber gets latest block number
func (n *node) LatestBlockNumber(ctx context.Context) (uint64, error) {
	var result uint64
	err := n.client.CallContext(ctx, &result, "kai_blockNumber")
	return result, err
}

func (n *node) getBlock(ctx context.Context, method string, args ...interface{}) (*Block, error) {
	var raw Block
	err := n.client.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, err
	}
	return &raw, nil
}

func (n *node) getBlockHeader(ctx context.Context, method string, args ...interface{}) (*Header, error) {
	var raw Header
	err := n.client.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, err
	}
	return &raw, nil
}

func (n *node) GetValidators(ctx context.Context, height uint64) (*types.ValidatorSet, error) {
	var valSet *types.ValidatorSet
	err := n.client.CallContext(ctx, &valSet, "kai_getValidatorSet", height)
	if err == nil && valSet == nil {
		err = kardia.NotFound
	}
	return valSet, err
}

func (n *node) GetCommit(ctx context.Context, height uint64) (*types.Commit, error) {
	var commit *types.Commit
	err := n.client.CallContext(ctx, &commit, "kai_getCommit", height)
	if err == nil && commit == nil {
		err = kardia.NotFound
	}
	return commit, err
}

func (n *node) GetProof(ctx context.Context, address common.Address, storageKeys []string, height uint64) (*kai.AccountResult, error) {
	var accountR *kai.AccountResult
	err := n.client.CallContext(ctx, &accountR, "kai_getProof", address, storageKeys, height, false)
	if err == nil && accountR == nil {
		err = kardia.NotFound
	}
	return accountR, err
}

type FullHeader struct {
	Header       *Header
	Commit       *types.Commit
	ValidatorSet *types.ValidatorSet
}

func (n *node) FullHeaderByNumber(ctx context.Context, height uint64) (*FullHeader, error) {
	header, err := n.BlockHeaderByNumber(ctx, height)
	if err != nil {
		return nil, err
	}
	validators, err := n.GetValidators(ctx, height)
	if err != nil {
		return nil, err
	}

	commit, err := n.GetCommit(ctx, height-1)
	if err != nil {
		return nil, err
	}
	return &FullHeader{
		Header:       header,
		ValidatorSet: validators,
		Commit:       commit,
	}, nil
}
