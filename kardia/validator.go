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

package kardia

import (
	"context"
	"math/big"
	"time"

	"github.com/kardiachain/go-kaiclient/types"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/mainchain/staking"

	"go.uber.org/zap"
)

// GetValidatorInfo returns information of this validator
func (ec *Client) GetValidatorInfo(ctx context.Context, valSmcAddr common.Address) (*staking.Validator, error) {
	payload, err := ec.validatorUtil.Abi.Pack("inforValidator")
	if err != nil {
		ec.lgr.Error("Error packing validator info payload: ", zap.Error(err))
		return nil, err
	}
	res, err := ec.KardiaCall(ctx, ec.contructCallArgs(valSmcAddr.Hex(), payload))
	if err != nil {
		ec.lgr.Error("GetValidatorInfo KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	var valInfo staking.Validator
	// unpack result
	err = ec.validatorUtil.Abi.UnpackIntoInterface(&valInfo, "inforValidator", res)
	if err != nil {
		ec.lgr.Error("Error unpacking validator info: ", zap.Error(err))
		return nil, err
	}
	return &valInfo, nil
}

// GetDelegationRewards returns reward of a delegation
func (ec *Client) GetDelegationRewards(ctx context.Context, valSmcAddr common.Address, delegatorAddr common.Address) (*big.Int, error) {
	payload, err := ec.validatorUtil.Abi.Pack("getDelegationRewards", delegatorAddr)
	if err != nil {
		ec.lgr.Error("Error packing delegation rewards payload: ", zap.Error(err))
		return nil, err
	}
	res, err := ec.KardiaCall(ctx, ec.contructCallArgs(valSmcAddr.Hex(), payload))
	if err != nil {
		ec.lgr.Error("GetDelegationRewards KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	var rewards struct {
		Rewards *big.Int
	}
	// unpack result
	err = ec.validatorUtil.Abi.UnpackIntoInterface(&rewards, "getDelegationRewards", res)
	if err != nil {
		ec.lgr.Error("Error unpacking delegation rewards: ", zap.Error(err))
		return nil, err
	}
	return rewards.Rewards, nil
}

// GetDelegatorStakedAmount returns staked amount of a delegator to current validator
func (ec *Client) GetDelegatorStakedAmount(ctx context.Context, valSmcAddr common.Address, delegatorAddr common.Address) (*big.Int, error) {
	payload, err := ec.validatorUtil.Abi.Pack("delegationByAddr", delegatorAddr)
	if err != nil {
		ec.lgr.Error("Error packing delegator staked amount payload: ", zap.Error(err))
		return nil, err
	}
	res, err := ec.KardiaCall(ctx, ec.contructCallArgs(valSmcAddr.Hex(), payload))
	if err != nil {
		ec.lgr.Error("GetDelegatorStakedAmount KardiaCall error: ", zap.Error(err))
		return nil, err
	}

	var delegation struct {
		Stake          *big.Int
		PreviousPeriod *big.Int
		Height         *big.Int
		Shares         *big.Int
		Owner          common.Address
	}
	// unpack result
	err = ec.validatorUtil.Abi.UnpackIntoInterface(&delegation, "delegationByAddr", res)
	if err != nil {
		ec.lgr.Error("Error unpacking delegator's staked amount: ", zap.Error(err))
		return nil, err
	}
	return delegation.Stake, nil
}

// GetUDBEntry returns unbonded amount and withdrawable amount of a delegation
func (ec *Client) GetUDBEntries(ctx context.Context, valSmcAddr common.Address, delegatorAddr common.Address) (*big.Int, *big.Int, error) {
	payload, err := ec.validatorUtil.Abi.Pack("getUBDEntries", delegatorAddr)
	if err != nil {
		ec.lgr.Error("Error packing UDB entry payload: ", zap.Error(err))
		return nil, nil, err
	}
	res, err := ec.KardiaCall(ctx, ec.contructCallArgs(valSmcAddr.Hex(), payload))
	if err != nil {
		ec.lgr.Error("GetUDBEntry KardiaCall error: ", zap.Error(err))
		return nil, nil, err
	}

	var udbEntry struct {
		Balances        []*big.Int
		CompletionTimes []*big.Int
	}
	// unpack result
	err = ec.validatorUtil.Abi.UnpackIntoInterface(&udbEntry, "getUBDEntries", res)
	if err != nil {
		ec.lgr.Error("Error unpacking UDB entry: ", zap.Error(err))
		return nil, nil, err
	}
	totalAmount := new(big.Int).SetInt64(0)
	withdrawableAmount := new(big.Int).SetInt64(0)
	now := new(big.Int).SetInt64(time.Now().Unix())
	for i, balance := range udbEntry.Balances {
		totalAmount = new(big.Int).Add(totalAmount, balance)
		if udbEntry.CompletionTimes[i].Cmp(now) == -1 {
			withdrawableAmount = new(big.Int).Add(withdrawableAmount, balance)
		}
	}
	if totalAmount == nil || len(totalAmount.Bits()) == 0 {
		totalAmount = new(big.Int).SetInt64(0)
	}
	if withdrawableAmount == nil || len(withdrawableAmount.Bits()) == 0 {
		withdrawableAmount = new(big.Int).SetInt64(0)
	}
	return totalAmount, withdrawableAmount, nil
}

func (ec *Client) GetValidatorParams(ctx context.Context, valSmcAddr common.Address) (*types.ValidatorParams, error) {
	payload, err := ec.validatorUtil.Abi.Pack("params")
	if err != nil {
		ec.lgr.Error("Error packing validator params payload: ", zap.Error(err))
		return nil, err
	}
	res, err := ec.KardiaCall(ctx, ec.contructCallArgs(valSmcAddr.Hex(), payload))
	if err != nil {
		ec.lgr.Error("GetDelegatorStakedAmount KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}

	var params types.ValidatorParams
	// unpack result
	err = ec.validatorUtil.Abi.UnpackIntoInterface(&params, "params", res)
	if err != nil {
		ec.lgr.Error("Error unpacking validator params: ", zap.Error(err))
		return nil, err
	}
	return &params, nil
}
