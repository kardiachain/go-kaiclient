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

	"github.com/kardiachain/go-kaiclient/types"
	"github.com/kardiachain/go-kardia/lib/common"
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
	valsSet, err := ec.GetValidatorSets(ctx)
	if err != nil {
		return nil, err
	}
	var valsList []*types.ValidatorsByDelegator
	for _, val := range valAddrs.ValAddrs {
		valInfo, err := ec.GetValidatorInfo(ctx, val)
		if err != nil {
			continue
		}
		var name []byte
		for _, b := range valInfo.Name {
			if b != 0 {
				name = append(name, byte(b))
			}
		}
		owner, err := ec.GetOwnerFromValidatorSMC(ctx, val)
		if err != nil {
			return nil, err
		}
		reward, err := ec.GetDelegationRewards(ctx, val, delAddr)
		if err != nil {
			continue
		}
		stakedAmount, err := ec.GetDelegatorStakedAmount(ctx, val, delAddr)
		if err != nil {
			continue
		}
		unbondedAmount, withdrawableAmount, err := ec.GetUDBEntries(ctx, val, delAddr)
		if err != nil {
			continue
		}
		// re-update validator role based on his status
		valInfo.Status, err = ec.getValidatorStatus(valsSet, valInfo)
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

// GetOwnerFromValidatorSMC returns owner address from validator contract address
func (ec *Client) GetOwnerFromValidatorSMC(ctx context.Context, valSmcAddr common.Address) (common.Address, error) {
	payload, err := ec.stakingUtil.Abi.Pack("valOf", valSmcAddr)
	if err != nil {
		ec.lgr.Error("Error packing owner of validator SMC payload: ", zap.Error(err))
		return common.Address{}, err
	}
	res, err := ec.KardiaCall(ctx, ec.contructCallArgs(ec.stakingUtil.ContractAddress.Hex(), payload))
	if err != nil {
		ec.lgr.Error("GetOwnerFromValidatorSMC KardiaCall error: ", zap.Error(err))
		return common.Address{}, err
	}
	if len(res) == 0 {
		ec.lgr.Debug("GetOwnerFromValidatorSMC KardiaCall empty result")
		return common.Address{}, ErrNotAValidatorAddress
	}
	var result struct {
		ValSmcAddr common.Address
	}
	err = ec.stakingUtil.Abi.UnpackIntoInterface(&result, "valOf", res)
	if err != nil {
		ec.lgr.Error("Error unpacking owner of validator SMC error: ", zap.Error(err))
		return common.Address{}, err
	}
	return result.ValSmcAddr, nil
}

// GetOwnerFromValidatorSMC returns owner address from validator contract address
func (ec *Client) GetValidatorSMCFromOwner(ctx context.Context, valAddr common.Address) (common.Address, error) {
	payload, err := ec.stakingUtil.Abi.Pack("ownerOf", valAddr)
	if err != nil {
		ec.lgr.Error("Error packing validator SMC of owner payload: ", zap.Error(err))
		return common.Address{}, err
	}
	res, err := ec.KardiaCall(ctx, ec.contructCallArgs(ec.stakingUtil.ContractAddress.Hex(), payload))
	if err != nil {
		ec.lgr.Error("GetValidatorSMCFromOwner KardiaCall error: ", zap.Error(err))
		return common.Address{}, err
	}
	if len(res) == 0 {
		ec.lgr.Debug("GetValidatorSMCFromOwner KardiaCall empty result")
		return common.Address{}, ErrNotAValidatorAddress
	}
	var result struct {
		ValSmcAddr common.Address
	}
	err = ec.stakingUtil.Abi.UnpackIntoInterface(&result, "ownerOf", res)
	if err != nil {
		ec.lgr.Error("Error unpacking validator SMC of owner error: ", zap.Error(err))
		return common.Address{}, err
	}
	return result.ValSmcAddr, nil
}

// GetValidatorSets returns current proposers set of network
func (ec *Client) GetValidatorSets(ctx context.Context) ([]common.Address, error) {
	payload, err := ec.stakingUtil.Abi.Pack("getValidatorSets")
	if err != nil {
		ec.lgr.Error("Error packing proposers list payload: ", zap.Error(err))
		return nil, err
	}
	res, err := ec.KardiaCall(ctx, ec.contructCallArgs(ec.stakingUtil.ContractAddress.Hex(), payload))
	if err != nil {
		ec.lgr.Error("GetValidatorSets KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	if len(res) == 0 {
		ec.lgr.Debug("GetValidatorSets KardiaCall empty result")
		return nil, nil
	}
	var result struct {
		ValAddrs []common.Address
		Powers   []*big.Int
	}
	err = ec.stakingUtil.Abi.UnpackIntoInterface(&result, "getValidatorSets", res)
	if err != nil {
		ec.lgr.Error("Error unpacking proposers list error: ", zap.Error(err))
		return nil, err
	}
	return result.ValAddrs, nil
}

// GetAllValsLength returns number of validators
func (ec *Client) GetAllValsLength(ctx context.Context) (*big.Int, error) {
	payload, err := ec.stakingUtil.Abi.Pack("allValsLength")
	if err != nil {
		ec.lgr.Error("Error packing get all validators length payload: ", zap.Error(err))
		return nil, err
	}

	res, err := ec.KardiaCall(ctx, ec.contructCallArgs(ec.stakingUtil.ContractAddress.Hex(), payload))
	if err != nil {
		ec.lgr.Error("GetAllValsLength KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	if len(res) == 0 {
		ec.lgr.Debug("GetAllValsLength KardiaCall empty result")
		return nil, ErrEmptyList
	}

	var valsLength *big.Int
	// unpack result
	err = ec.stakingUtil.Abi.UnpackIntoInterface(&valsLength, "allValsLength", res)
	if err != nil {
		ec.lgr.Error("Error unpacking get all validators length error: ", zap.Error(err))
		return nil, err
	}
	return valsLength, nil
}

// GetValSmcAddr returns validator's info based on his index
func (ec *Client) GetValSmcAddr(ctx context.Context, index *big.Int) (common.Address, error) {
	payload, err := ec.stakingUtil.Abi.Pack("allVals", index)
	if err != nil {
		ec.lgr.Error("Error packing get validator SMC address payload: ", zap.Error(err))
		return common.Address{}, err
	}
	res, err := ec.KardiaCall(ctx, ec.contructCallArgs(ec.stakingUtil.ContractAddress.Hex(), payload))
	if err != nil {
		ec.lgr.Error("GetValSmcAddr KardiaCall error: ", zap.Error(err))
		return common.Address{}, err
	}
	if len(res) == 0 {
		ec.lgr.Debug("GetOwnerFromValidatorSMC KardiaCall empty result")
		return common.Address{}, nil
	}

	var valSmc struct {
		AddrValSmc common.Address
	}

	err = ec.stakingUtil.Abi.UnpackIntoInterface(&valSmc, "allVals", res)
	if err != nil {
		ec.lgr.Error("Error unpacking get validator SMC address error: ", zap.Error(err))
		return common.Address{}, err
	}

	return valSmc.AddrValSmc, nil
}

// GetValFromOwner returns address validator smc of validator
func (ec *Client) GetValFromOwner(ctx context.Context, valAddr common.Address) (common.Address, error) {
	payload, err := ec.stakingUtil.Abi.Pack("ownerOf", valAddr)
	if err != nil {
		ec.lgr.Error("Error packing get validator SMC address from owner payload: ", zap.Error(err))
		return common.Address{}, err
	}
	res, err := ec.KardiaCall(ctx, ec.contructCallArgs(ec.stakingUtil.ContractAddress.Hex(), payload))
	if err != nil {
		ec.lgr.Error("GetValFromOwner KardiaCall error: ", zap.Error(err))
		return common.Address{}, err
	}
	if len(res) == 0 {
		ec.lgr.Debug("GetValFromOwner KardiaCall empty result")
		return common.Address{}, nil
	}

	var valSmc struct {
		AddrValSmc common.Address
	}
	err = ec.stakingUtil.Abi.UnpackIntoInterface(&valSmc, "ownerOf", res)
	if err != nil {
		ec.lgr.Error("Error unpacking get validator SMC address from owner error: ", zap.Error(err))
		return common.Address{}, err
	}

	return valSmc.AddrValSmc, nil
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
