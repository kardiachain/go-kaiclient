// Package kardia
package kardia

import (
	"context"
	"math/big"
	"time"

	"github.com/kardiachain/go-kardia/lib/common"
	"go.uber.org/zap"
)

type IDelegator interface {
	UnbondedRecords(ctx context.Context, validatorSMCAddress common.Address, delegatorAddress common.Address) (*big.Int, *big.Int, error)
}

func (n *node) UnbondedRecords(ctx context.Context, validatorSMCAddress common.Address, delegatorAddress common.Address) (*big.Int, *big.Int, error) {
	payload, err := n.validatorSMC.Abi.Pack("getUBDEntries", delegatorAddress)
	if err != nil {
		n.lgr.Error("Error packing UDB entry payload: ", zap.Error(err))
		return nil, nil, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(validatorSMCAddress.Hex(), payload))
	if err != nil {
		n.lgr.Error("GetUDBEntry KardiaCall error: ", zap.Error(err))
		return nil, nil, err
	}
	if len(res) == 0 {
		return nil, nil, ErrEmptyList
	}

	var result struct {
		Balances        []*big.Int
		CompletionTimes []*big.Int
	}
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&result, "getUBDEntries", res)
	if err != nil {
		n.lgr.Error("Error unpacking UDB entry: ", zap.Error(err))
		return nil, nil, err
	}
	totalAmount := new(big.Int).SetInt64(0)
	withdrawableAmount := new(big.Int).SetInt64(0)
	now := new(big.Int).SetInt64(time.Now().Unix())
	for i, balance := range result.Balances {
		if result.CompletionTimes[i].Cmp(now) == -1 {
			withdrawableAmount = new(big.Int).Add(withdrawableAmount, balance)
		} else {
			totalAmount = new(big.Int).Add(totalAmount, balance)
		}
	}
	return totalAmount, withdrawableAmount, nil
}
