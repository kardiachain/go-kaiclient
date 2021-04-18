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
	"fmt"
	"math/big"
	"time"

	"github.com/kardiachain/go-kardia/lib/common"
	"go.uber.org/zap"
)

type IValidator interface {
	ValidatorInfo(ctx context.Context, validatorSMCAddress string) (*Validator, error)
	SigningInfo(ctx context.Context, validatorSMCAddress string) (*SigningInfo, error)
	DelegatorAddresses(ctx context.Context, validatorSMCAddress string) ([]common.Address, error)
	DelegatorsWithShare(ctx context.Context, validatorSMCAddress string) ([]*DelegatorWithShare, error)
	DelegationRewards(ctx context.Context, validatorSMCAddr, delegatorAddress string) (*big.Int, error)
	DelegatorStakedAmount(ctx context.Context, validatorSMCAddress, delegatorAddress string) (*big.Int, error)
	ValidatorCommission(ctx context.Context, valSmcAddr string) (*Commission, error)
	SlashEvents(ctx context.Context, validatorSMCAddress string) ([]*SlashEvents, error)

	//deprecated
	Validators(ctx context.Context) ([]*Validator, error)
	Validator(ctx context.Context, validatorSMCAddress string) (*Validator, error)
	ValidatorSets(ctx context.Context) ([]common.Address, error)
	SMCAddressOfValidator(ctx context.Context, validatorAddress string) (common.Address, error)
	ValidatorAddressOfSMC(ctx context.Context, validatorSMCAddress string) (common.Address, error)
}

func (n *node) ValidatorInfo(ctx context.Context, validatorSMCAddress string) (*Validator, error) {
	lgr := n.lgr.With(zap.String("method", "Validator"))
	startLoadInfo := time.Now()
	payload, err := n.validatorSMC.Abi.Pack("inforValidator")
	if err != nil {
		lgr.Error("Error packing validator info payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(validatorSMCAddress, payload))
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
	lgr.Debug("Finished load validator info", zap.Duration("Time", time.Now().Sub(startLoadInfo)))
	return &valInfo, nil
}

//SigningInfo
func (n *node) SigningInfo(ctx context.Context, validatorSMCAddress string) (*SigningInfo, error) {
	lgr := n.lgr.With(zap.String("method", "getSigningInfo"))
	payload, err := n.validatorSMC.Abi.Pack("signingInfo")
	if err != nil {
		lgr.Error("Error packing get signingInfo payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(validatorSMCAddress, payload))
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

func (n *node) DelegationRewards(ctx context.Context, validatorSMCAddr, delegatorAddress string) (*big.Int, error) {
	payload, err := n.validatorSMC.Abi.Pack("getDelegationRewards", common.HexToAddress(delegatorAddress))
	if err != nil {
		n.lgr.Error("Error packing delegation rewards payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(validatorSMCAddr, payload))
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

// SlashEvents returns detailed all slash events of this validator
func (n *node) SlashEvents(ctx context.Context, validatorSMCAddress string) ([]*SlashEvents, error) {
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
		res, err := n.KardiaCall(ctx, ConstructCallArgs(validatorSMCAddress, payload))
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

// DelegatorStakedAmount returns staked amount of a delegator to current validator
func (n *node) DelegatorStakedAmount(ctx context.Context, validatorSMCAddress, delegatorAddress string) (*big.Int, error) {
	payload, err := n.validatorSMC.Abi.Pack("delegationByAddr", common.HexToAddress(delegatorAddress))
	if err != nil {
		n.lgr.Error("Error packing delegator staked amount payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(validatorSMCAddress, payload))
	if err != nil {
		n.lgr.Error("DelegatorStakedAmount KardiaCall error: ", zap.Error(err))
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
func (n *node) ValidatorCommission(ctx context.Context, valSmcAddr string) (*Commission, error) {
	payload, err := n.validatorSMC.Abi.Pack("commission")
	if err != nil {
		return nil, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(valSmcAddr, payload))
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

func (n *node) DelegatorAddresses(ctx context.Context, validatorSMCAddress string) ([]common.Address, error) {
	payload, err := n.validatorSMC.Abi.Pack("getDelegations")
	if err != nil {
		return nil, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(validatorSMCAddress, payload))
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
	return result.Addresses, nil
}

func (n *node) DelegatorsWithShare(ctx context.Context, validatorSMCAddress string) ([]*DelegatorWithShare, error) {
	payload, err := n.validatorSMC.Abi.Pack("getDelegations")
	if err != nil {
		return nil, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(validatorSMCAddress, payload))
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
	var delegatorsWithShare []*DelegatorWithShare
	for id := range result.Addresses {
		delegatorsWithShare = append(delegatorsWithShare,
			&DelegatorWithShare{Address: result.Addresses[id], Share: result.Shares[id]},
		)
	}
	return delegatorsWithShare, nil
}

// getSlashEventsSize returns number of slash events of this validator
func (n *node) getSlashEventsSize(ctx context.Context, validatorSMCAddress string) (int, error) {
	payload, err := n.validatorSMC.Abi.Pack("getSlashEventsLength")
	if err != nil {
		n.lgr.Error("Error packing get slash events length payload: ", zap.Error(err))
		return 0, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(validatorSMCAddress, payload))
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

func (n *node) ValidatorSets(ctx context.Context) ([]common.Address, error) {
	payload, err := n.stakingSMC.Abi.Pack("getValidatorSets")
	if err != nil {
		n.lgr.Error("Error packing proposers list payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(n.stakingSMC.ContractAddress.Hex(), payload))
	if err != nil {
		n.lgr.Error("GetValidatorSets KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	if len(res) == 0 {
		n.lgr.Debug("GetValidatorSets KardiaCall empty result")
		return nil, nil
	}
	var result struct {
		ValidatorAddresses []common.Address
		ValidatorPowers    []*big.Int
	}
	err = n.stakingSMC.Abi.UnpackIntoInterface(&result, "getValidatorSets", res)
	if err != nil {
		n.lgr.Error("Error unpacking proposers list error: ", zap.Error(err))
		return nil, err
	}
	return result.ValidatorAddresses, nil
}

//deprecated
func (n *node) Validators(ctx context.Context) ([]*Validator, error) {
	lgr := n.lgr.With(zap.String("method", "Validators"))
	var (
		validators []*Validator
	)

	validatorSMCAddresses, err := n.ValidatorSMCAddresses(ctx)
	if err != nil {
		return nil, err
	}

	for _, smcAddr := range validatorSMCAddresses {
		loadValidatorStartTime := time.Now()
		v, err := n.Validator(ctx, smcAddr.Hex())
		if err != nil {
			return nil, err
		}
		lgr.Debug("Finished load validator", zap.String("Validator", fmt.Sprintf("%s", v.Name)), zap.String("SMCAddress", v.SMCAddress.String()), zap.Int("TotalDelegator", len(v.Delegators)), zap.Duration("Total", time.Now().Sub(loadValidatorStartTime)))
		validators = append(validators, v)
	}
	return validators, nil
}

// deprecated
func (n *node) Validator(ctx context.Context, validatorSMCAddress string) (*Validator, error) {
	lgr := n.lgr.With(zap.String("method", "Validator"))
	startLoadInfo := time.Now()
	payload, err := n.validatorSMC.Abi.Pack("inforValidator")
	if err != nil {
		lgr.Error("Error packing validator info payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(validatorSMCAddress, payload))
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
	lgr.Debug("Finished load validator info", zap.Duration("Total", time.Now().Sub(startLoadInfo)))
	valInfo.SMCAddress = common.HexToAddress(validatorSMCAddress)
	startGetValidatorComm := time.Now()
	commission, err := n.ValidatorCommission(ctx, validatorSMCAddress)
	if err != nil {
		return nil, err
	}
	valInfo.Commission = commission
	lgr.Debug("Finished load get validator comm", zap.Duration("Total", time.Now().Sub(startGetValidatorComm)))

	startGetSigningInfo := time.Now()
	signingInfo, err := n.SigningInfo(ctx, validatorSMCAddress)
	if err != nil {
		return nil, err
	}
	valInfo.SigningInfo = signingInfo
	lgr.Debug("Finished load get signing info", zap.Duration("Time", time.Now().Sub(startGetSigningInfo)))
	return &valInfo, nil
}

// deprecated
func (n *node) delegatorsOfValidator(ctx context.Context, validatorSMCAddress string) ([]*Delegator, error) {
	lgr := n.lgr.With(zap.String("method", "delegatorsOfValidator"))
	payload, err := n.validatorSMC.Abi.Pack("getDelegations")
	if err != nil {
		return nil, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(validatorSMCAddress, payload))
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
		startGetReward := time.Now()
		delegatorAddress := delAddr.Hex()
		reward, err := n.DelegationRewards(ctx, validatorSMCAddress, delegatorAddress)
		if err != nil {
			continue
		}
		lgr.Debug("Finished load reward", zap.Duration("Total", time.Now().Sub(startGetReward)))

		startGetStakedAmount := time.Now()
		stakedAmount, err := n.DelegatorStakedAmount(ctx, validatorSMCAddress, delegatorAddress)
		if err != nil {
			continue
		}
		lgr.Debug("Finished load staked amount", zap.Duration("Total", time.Now().Sub(startGetStakedAmount)))
		delegators = append(delegators, &Delegator{
			Address:      delAddr,
			StakedAmount: stakedAmount,
			Reward:       reward,
		})
	}
	return delegators, nil
}

func (n *node) Delegator(ctx context.Context, validatorSMCAddress, delegatorAddress common.Address) (*Delegator, error) {
	lgr := n.lgr.With(zap.String("method", "Delegator"))
	startGetReward := time.Now()
	reward, err := n.DelegationRewards(ctx, validatorSMCAddress.Hex(), delegatorAddress.Hex())
	if err != nil {
		return nil, err
	}
	lgr.Debug("Finished load reward", zap.Duration("Total", time.Now().Sub(startGetReward)))

	startGetStakedAmount := time.Now()
	stakedAmount, err := n.DelegatorStakedAmount(ctx, validatorSMCAddress.Hex(), delegatorAddress.Hex())
	if err != nil {
		return nil, err
	}
	lgr.Debug("Finished load staked amount", zap.Duration("Total", time.Now().Sub(startGetStakedAmount)))
	d := &Delegator{
		Address:      delegatorAddress,
		StakedAmount: stakedAmount,
		Reward:       reward,
	}
	return d, nil
}

func (n *node) SMCAddressOfValidator(ctx context.Context, validatorAddress string) (common.Address, error) {
	lgr := n.lgr.With(zap.String("method", "SMCAddressOfValidator"))

	payload, err := n.stakingSMC.Abi.Pack("ownerOf", common.HexToAddress(validatorAddress))
	if err != nil {
		lgr.Error("Error packing validator SMC of owner payload: ", zap.Error(err))
		return common.Address{}, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(n.stakingSMC.ContractAddress.Hex(), payload))
	if err != nil {
		lgr.Error("GetValidatorSMCFromOwner KardiaCall error: ", zap.Error(err))
		return common.Address{}, err
	}

	if len(res) == 0 {
		lgr.Debug("GetValidatorSMCFromOwner KardiaCall empty result")
		return common.Address{}, fmt.Errorf("found nothing")
	}
	var result struct {
		ValSmcAddr common.Address
	}
	err = n.stakingSMC.Abi.UnpackIntoInterface(&result, "ownerOf", res)
	if err != nil {
		lgr.Error("Error unpacking validator SMC of owner error: ", zap.Error(err))
		return common.Address{}, err
	}
	return result.ValSmcAddr, nil
}

func (n *node) ValidatorAddressOfSMC(ctx context.Context, validatorSMCAddress string) (common.Address, error) {
	lgr := n.lgr.With(zap.String("method", "ValidatorAddressOfSMC"))
	payload, err := n.stakingSMC.Abi.Pack("valOf", common.HexToAddress(validatorSMCAddress))
	if err != nil {
		lgr.Error("Error packing owner of validator SMC payload: ", zap.Error(err))
		return common.Address{}, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(n.stakingSMC.ContractAddress.Hex(), payload))
	if err != nil {
		lgr.Error("GetOwnerFromValidatorSMC KardiaCall error: ", zap.Error(err))
		return common.Address{}, err
	}
	if len(res) == 0 {
		lgr.Debug("GetOwnerFromValidatorSMC KardiaCall empty result")
		return common.Address{}, fmt.Errorf("cannot found")
	}
	var result struct {
		ValSmcAddr common.Address
	}
	err = n.stakingSMC.Abi.UnpackIntoInterface(&result, "valOf", res)
	if err != nil {
		lgr.Error("Error unpacking owner of validator SMC error: ", zap.Error(err))
		return common.Address{}, err
	}
	return result.ValSmcAddr, nil
}
