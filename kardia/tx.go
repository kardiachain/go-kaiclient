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

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/kardiachain/go-kardia"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/rlp"
	"github.com/kardiachain/go-kardia/types"
)

type ITx interface {
	GetTransaction(ctx context.Context, hash string) (*Transaction, error)
	GetTransactionReceipt(ctx context.Context, txHash string) (*Receipt, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	SendRawTransaction(ctx context.Context, tx *types.Transaction) error
}

// GetTransaction returns the transaction with the given hash.
func (n *node) GetTransaction(ctx context.Context, hash string) (*Transaction, error) {
	var raw *Transaction
	err := n.client.CallContext(ctx, &raw, "tx_getTransaction", common.HexToHash(hash))
	if err != nil {
		return nil, err
	} else if raw == nil {
		return nil, kardia.NotFound
	}
	return raw, nil
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (n *node) GetTransactionReceipt(ctx context.Context, txHash string) (*Receipt, error) {
	var r *Receipt
	err := n.client.CallContext(ctx, &r, "tx_getTransactionReceipt", common.HexToHash(txHash))
	if err == nil {
		if r == nil {
			return nil, kardia.NotFound
		}
	}
	return r, err
}

// SendRawTransaction injects a signed transaction into the pending pool for execution.
// If the transaction was a contract creation use the GetTransactionReceipt method to get the
// contract address after the transaction has been mined.
func (n *node) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}
	return n.client.CallContext(ctx, nil, "tx_sendRawTransaction", hexutil.Encode(data))
}

// SendRawTransaction injects a signed transaction into the pending pool for execution.
// If the transaction was a contract creation use the GetTransactionReceipt method to get the
// contract address after the transaction has been mined.
func (n *node) SendRawTransaction(ctx context.Context, tx *types.Transaction) error {
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}
	return n.client.CallContext(ctx, nil, "tx_sendRawTransaction", hexutil.Encode(data))
}
