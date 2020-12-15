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

	"go.uber.org/zap"

	"github.com/kardiachain/go-kardia/lib/common"

	"github.com/kardiachain/go-kaiclient/types"
)

func (ec *Client) GetValidatorsByDelegator(ctx context.Context, delAddr common.Address) ([]*types.ValidatorsByDelegator, error) {
	// construct input data
	payload, err := ec.stakingUtil.Abi.Pack("getValidatorsByDelegator", delAddr)
	if err != nil {
		return nil, err
	}
	// make static call through `kai_kardiaCall` API
	res, err := ec.KardiaCall(ctx, ec.contructCallArgs(ec.stakingUtil.ContractAddress.Hex(), payload))
	if err != nil {
		return nil, err
	}
	// get validators list of delegator
	var valAddrs struct {
		ValAddrs []common.Address
	}
	// unpacking data
	err = ec.stakingUtil.Abi.UnpackIntoInterface(&valAddrs, "getValidatorsByDelegator", res)
	if err != nil {
		return nil, err
	}

	// gather additional information about validators
	var valsList []*types.ValidatorsByDelegator
	for _, val := range valAddrs.ValAddrs {
		valInfo, err := ec.GetValidatorInfo(ctx, val)
		if err != nil {
			return nil, err
		}
		var name []byte
		for _, b := range valInfo.Name {
			if b != 0 {
				name = append(name, byte(b))
			}
		}
		owner, err := ec.GetValidatorContractFromOwner(ctx, val)
		if err != nil {
			return nil, err
		}
		reward, err := ec.GetDelegationRewards(ctx, val, delAddr)
		if err != nil {
			return nil, err
		}
		stakedAmount, err := ec.GetDelegatorStakedAmount(ctx, val, delAddr)
		if err != nil {
			return nil, err
		}
		unbondedAmount, withdrawableAmount, err := ec.GetUDBEntries(ctx, val, delAddr)
		if err != nil {
			return nil, err
		}
		validator := &types.ValidatorsByDelegator{
			Name:                  string(name),
			Validator:             owner,
			ValidatorContractAddr: val,
			ValidatorStatus:       valInfo.Status,
			StakedAmount:          stakedAmount.String(),
			ClaimableRewards:      reward.String(),
			UnbondedAmount:        unbondedAmount.String(),
			WithdrawableAmount:    withdrawableAmount.String(),
		}
		valsList = append(valsList, validator)
	}
	return valsList, nil
}

// GetValidatorContractFromOwner returns validator contract address from owner address
func (ec *Client) GetValidatorContractFromOwner(ctx context.Context, valAddr common.Address) (common.Address, error) {
	payload, err := ec.stakingUtil.Abi.Pack("valOf", valAddr)
	if err != nil {
		ec.lgr.Error("Error packing owner of validator SMC payload: ", zap.Error(err))
		return common.Address{}, err
	}
	res, err := ec.KardiaCall(ctx, ec.contructCallArgs(ec.stakingUtil.ContractAddress.Hex(), payload))
	if err != nil {
		ec.lgr.Error("GetDelegationRewards KardiaCall error: ", zap.Error(err))
		return common.Address{}, err
	}
	var owner struct {
		ValSmcAddr common.Address
	}
	err = ec.stakingUtil.Abi.UnpackIntoInterface(&owner, "valOf", res)
	if err != nil {
		ec.lgr.Error("Error unpacking owner of validator SMC error: ", zap.Error(err))
		return common.Address{}, err
	}
	return owner.ValSmcAddr, nil
}

func (ec *Client) contructCallArgs(address string, payload []byte) types.CallArgsJSON {
	return types.CallArgsJSON{
		From:     address,
		To:       &address,
		Gas:      100000000,
		GasPrice: big.NewInt(0),
		Value:    big.NewInt(0),
		Data:     common.Bytes(payload).String(),
	}
}
