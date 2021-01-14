// Package kardia
package kardia

import (
	"context"

	"github.com/kardiachain/go-kardia/lib/common"
)

// BlockByHash returns the given full block.
//
// Use HeaderByHash if you don't need all transactions or uncle headers.
func (n *node) BlockByHash(ctx context.Context, hash string) (*Block, error) {
	return n.getBlock(ctx, "kai_getBlockByHash", common.HexToHash(hash))
}

// BlockByHeight returns a block from the current canonical chain.
//
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
