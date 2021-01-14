// Package kardia
package kardia

import (
	"context"

	"github.com/kardiachain/go-kardia"
	"github.com/kardiachain/go-kardia/lib/common"
)

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
//
// If the transaction was a contract creation use the GetTransactionReceipt method to get the
// contract address after the transaction has been mined.
func (n *node) SendRawTransaction(ctx context.Context, tx string) error {
	return n.client.CallContext(ctx, nil, "tx_sendRawTransaction", tx)
}
