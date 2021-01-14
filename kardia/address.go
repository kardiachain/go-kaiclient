// Package kardia
package kardia

import (
	"context"

	"github.com/kardiachain/go-kardia/lib/common"
)

type IAddress interface {
	Balance(ctx context.Context, addressHash string) (string, error)
	StorageAt(ctx context.Context, addressHash string, key string) ([]byte, error)
	Code(ctx context.Context, addressHash string) (common.Bytes, error)
	NonceAt(ctx context.Context, addressHash string) (uint64, error)
}

func (n *node) Balance(ctx context.Context, addressHash string) (string, error) {
	var (
		result string
		err    error
	)
	err = n.client.CallContext(ctx, &result, "account_balance", common.HexToAddress(addressHash), "latest")
	return result, err
}

// StorageAt returns the value of key in the contract storage of the given account.
// The block number can be nil, in which case the value is taken from the latest known block.
func (n *node) StorageAt(ctx context.Context, addressHash string, key string) ([]byte, error) {
	var result common.Bytes
	err := n.client.CallContext(ctx, &result, "account_getStorageAt", common.HexToAddress(addressHash), key, "latest")
	return result, err
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
func (n *node) Code(ctx context.Context, addressHash string) (common.Bytes, error) {
	var result common.Bytes
	err := n.client.CallContext(ctx, &result, "account_getCode", common.HexToAddress(addressHash), "latest")
	return result, err
}

// NonceAt returns the account nonce of the given account.
func (n *node) NonceAt(ctx context.Context, account string) (uint64, error) {
	var result uint64
	err := n.client.CallContext(ctx, &result, "account_nonce", common.HexToAddress(account))
	return result, err
}
