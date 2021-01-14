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

	"github.com/kardiachain/go-kardia/lib/common"
	"go.uber.org/zap"
)

type IValidator interface {
	Validator(ctx context.Context, validatorSMCAddr string) (*Validator, error)
	Validators(ctx context.Context) ([]*Validator, error)
}

func (n *node) Validators(ctx context.Context) ([]*Validator, error) {
	var (
		validators []*Validator
	)

	validatorSMCAddresses, err := n.validatorSMCAddresses(ctx)
	if err != nil {
		return nil, err
	}

	for _, smcAddr := range validatorSMCAddresses {
		v, err := n.Validator(ctx, smcAddr.Hex())
		if err != nil {
			return nil, err
		}
		validators = append(validators, v)
	}
	return validators, nil
}

func (n *node) Validator(ctx context.Context, validatorSMCAddress string) (*Validator, error) {
	lgr := n.lgr.With(zap.String("method", "Validator"))
	payload, err := n.validatorSMC.Abi.Pack("inforValidator")
	if err != nil {
		lgr.Error("Error packing validator info payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.KardiaCall(ctx, constructCallArgs(validatorSMCAddress, payload))
	if err != nil {
		lgr.Error("GetValidatorInfo KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	var valInfo Validator
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&valInfo, "inforValidator", res)
	if err != nil {
		lgr.Error("Error unpacking validator info: ", zap.Error(err))
		return nil, err
	}

	valInfo.SMCAddress = common.HexToAddress(validatorSMCAddress)

	commission, err := n.getValidatorCommission(ctx, validatorSMCAddress)
	if err != nil {
		return nil, err
	}
	valInfo.Commission = commission

	signingInfo, err := n.getSigningInfo(ctx, validatorSMCAddress)
	if err != nil {
		return nil, err
	}
	valInfo.SigningInfo = signingInfo

	delegators, err := n.getDelegators(ctx, validatorSMCAddress)
	if err != nil {
		return nil, err
	}
	valInfo.Delegators = delegators

	return &valInfo, nil
}

// Helper
func (n *node) getSigningInfo(ctx context.Context, validatorSMCAddress string) (*SigningInfo, error) {
	lgr := n.lgr.With(zap.String("method", "getSigningInfo"))
	payload, err := n.validatorSMC.Abi.Pack("signingInfo")
	if err != nil {
		lgr.Error("Error packing get signingInfo payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.KardiaCall(ctx, constructCallArgs(validatorSMCAddress, payload))
	if err != nil {
		lgr.Error("GetSigningInfo KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	var result SigningInfo
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&result, "signingInfo", res)
	if err != nil {
		lgr.Error("Error unpack get signingInfo: ", zap.Error(err))
		return nil, err
	}
	return &result, nil
}

// getDelegatorStakedAmount returns staked amount of a delegator to current validator
func (n *node) getDelegatorStakedAmount(ctx context.Context, valSmcAddr, delegatorAddress string) (*big.Int, error) {
	payload, err := n.validatorSMC.Abi.Pack("delegationByAddr", common.HexToAddress(delegatorAddress))
	if err != nil {
		n.lgr.Error("Error packing delegator staked amount payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.KardiaCall(ctx, constructCallArgs(valSmcAddr, payload))
	if err != nil {
		n.lgr.Error("getDelegatorStakedAmount KardiaCall error: ", zap.Error(err))
		return nil, err
	}

	var result struct {
		Stake          *big.Int
		PreviousPeriod *big.Int
		Height         *big.Int
		Shares         *big.Int
		Owner          common.Address
	}
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&result, "delegationByAddr", res)
	if err != nil {
		n.lgr.Error("Error unpacking delegator's staked amount: ", zap.Error(err))
		return nil, err
	}
	return result.Stake, nil
}

// GetValidator show info of a validator based on address
func (n *node) getValidatorCommission(ctx context.Context, valSmcAddr string) (*Commission, error) {
	payload, err := n.validatorSMC.Abi.Pack("commission")
	if err != nil {
		return nil, err
	}
	res, err := n.KardiaCall(ctx, constructCallArgs(valSmcAddr, payload))
	if err != nil {
		return nil, err
	}

	var result Commission
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&result, "commission", res)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (n *node) getDelegators(ctx context.Context, validatorSMCAddress string) ([]*Delegator, error) {
	payload, err := n.validatorSMC.Abi.Pack("getDelegations")
	if err != nil {
		return nil, err
	}
	res, err := n.KardiaCall(ctx, constructCallArgs(validatorSMCAddress, payload))
	if err != nil {
		return nil, err
	}

	var result struct {
		Addresses []common.Address
		Shares    []*big.Int
	}
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&result, "getDelegations", res)
	if err != nil {
		n.lgr.Error("Error unpacking delegation details", zap.Error(err))
		return nil, err
	}
	var delegators []*Delegator
	for _, delAddr := range result.Addresses {
		delegatorAddress := delAddr.Hex()
		reward, err := n.getDelegationRewards(ctx, validatorSMCAddress, delegatorAddress)
		if err != nil {
			continue
		}
		stakedAmount, err := n.getDelegatorStakedAmount(ctx, validatorSMCAddress, delegatorAddress)
		if err != nil {
			continue
		}
		delegators = append(delegators, &Delegator{
			Address:      delAddr,
			StakedAmount: stakedAmount,
			Reward:       reward,
		})
	}
	return delegators, nil
}

func (n *node) getDelegationRewards(ctx context.Context, validatorSMCAddr, delegatorAddress string) (*big.Int, error) {
	payload, err := n.validatorSMC.Abi.Pack("getDelegationRewards", common.HexToAddress(delegatorAddress))
	if err != nil {
		n.lgr.Error("Error packing delegation rewards payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.KardiaCall(ctx, constructCallArgs(validatorSMCAddr, payload))
	if err != nil {
		n.lgr.Error("GetDelegationRewards KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	var result struct {
		Rewards *big.Int
	}
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&result, "getDelegationRewards", res)
	if err != nil {
		n.lgr.Error("Error unpacking delegation rewards: ", zap.Error(err))
		return nil, err
	}
	return result.Rewards, nil
}

// GetSlashEvents returns detailed all slash events of this validator
func (n *node) GetSlashEvents(ctx context.Context, validatorSMCAddress string) ([]*SlashEvents, error) {
	var events []*SlashEvents
	eventsSize, err := n.getSlashEventsSize(ctx, validatorSMCAddress)
	if err != nil {
		return nil, err
	}
	for i := 0; i < eventsSize; i++ {
		payload, err := n.validatorSMC.Abi.Pack("slashEvents", i)
		if err != nil {
			return nil, err
		}
		res, err := n.KardiaCall(ctx, constructCallArgs(validatorSMCAddress, payload))
		if err != nil {
			return nil, err
		}
		var result struct {
			Period   *big.Int
			Fraction *big.Int
			Height   *big.Int
		}
		// unpack result
		err = n.validatorSMC.Abi.UnpackIntoInterface(&result, "slashEvents", res)
		if err != nil {
			n.lgr.Error("Error unpacking slash event", zap.Error(err))
			return nil, err
		}
		events = append(events, &SlashEvents{
			Period:   result.Period.String(),
			Fraction: result.Fraction.String(),
			Height:   result.Height.String(),
		})
	}
	return events, nil
}

// getSlashEventsSize returns number of slash events of this validator
func (n *node) getSlashEventsSize(ctx context.Context, validatorSMCAddress string) (int, error) {
	payload, err := n.validatorSMC.Abi.Pack("getSlashEventsLength")
	if err != nil {
		n.lgr.Error("Error packing get slash events length payload: ", zap.Error(err))
		return 0, err
	}
	res, err := n.KardiaCall(ctx, constructCallArgs(validatorSMCAddress, payload))
	if err != nil {
		n.lgr.Warn("GetSlashEventsLength KardiaCall error: ", zap.Error(err))
		return 0, err
	}
	if len(res) == 0 {
		n.lgr.Debug("GetSlashEventsLength KardiaCall empty result")
		return 0, ErrEmptyList
	}

	var slashEventsSize int
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&slashEventsSize, "getSlashEventsLength", res)
	if err != nil {
		n.lgr.Error("Error unpacking get slash events length error: ", zap.Error(err))
		return 0, err
	}
	return slashEventsSize, nil
}

func (n *node) getValidatorSets(ctx context.Context) ([]common.Address, error) {
	payload, err := n.stakingSMC.Abi.Pack("getValidatorSets")
	if err != nil {
		n.lgr.Error("Error packing proposers list payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.KardiaCall(ctx, constructCallArgs(n.stakingSMC.ContractAddress.Hex(), payload))
	if err != nil {
		n.lgr.Error("GetValidatorSets KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	if len(res) == 0 {
		n.lgr.Debug("GetValidatorSets KardiaCall empty result")
		return nil, nil
	}
	var result struct {
		ValAddrs []common.Address
		Powers   []*big.Int
	}
	err = n.stakingSMC.Abi.UnpackIntoInterface(&result, "getValidatorSets", res)
	if err != nil {
		n.lgr.Error("Error unpacking proposers list error: ", zap.Error(err))
		return nil, err
	}
	return result.ValAddrs, nil
}
