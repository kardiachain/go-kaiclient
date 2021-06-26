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

type IStaking interface {
	TotalStakedAmount(ctx context.Context) (*big.Int, error)
	ValidatorSMCAddresses(ctx context.Context) ([]common.Address, error)
}

func (n *node) TotalStakedAmount(ctx context.Context) (*big.Int, error) {
	payload, err := n.stakingSMC.Abi.Pack("totalBonded")
	if err != nil {
		n.lgr.Error("Error packing UDB entry payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(n.stakingSMC.ContractAddress.Hex(), payload))
	if err != nil {
		n.lgr.Error("Get totalBonded KardiaCall error: ", zap.Error(err))
		return nil, err
	}

	var result struct {
		TotalBonded *big.Int
	}
	// unpack result
	err = n.stakingSMC.Abi.UnpackIntoInterface(&result, "totalBonded", res)
	if err != nil {
		n.lgr.Error("Error unpacking UDB entry: ", zap.Error(err))
		return nil, err
	}
	return result.TotalBonded, nil
}

func (n *node) ValidatorSMCAddresses(ctx context.Context) ([]common.Address, error) {
	payload, err := n.stakingSMC.Abi.Pack("getAllValidator")
	if err != nil {
		return nil, err
	}

	res, err := n.KardiaCall(ctx, ConstructCallArgs(n.stakingSMC.ContractAddress.Hex(), payload))
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, ErrEmptyList
	}

	var validatorSMCAddresses []common.Address
	// unpack result
	err = n.stakingSMC.Abi.UnpackIntoInterface(&validatorSMCAddresses, "getAllValidator", res)
	if err != nil {
		return nil, err
	}

	return validatorSMCAddresses, nil
}

func (n *node) GetCirculatingSupply(ctx context.Context) (*big.Int, error) {
	payload, err := n.stakingSMC.Abi.Pack("totalSupply")
	if err != nil {
		n.lgr.Error("Error packing circulating supply payload: ", zap.Error(err))
		return nil, err
	}

	res, err := n.KardiaCall(ctx, ConstructCallArgs(n.stakingSMC.ContractAddress.Hex(), payload))
	if err != nil {
		n.lgr.Error("GetCirculatingSupply KardiaCall error: ", zap.Error(err))
		return nil, err
	}

	var result struct {
		TotalSupply *big.Int
	}
	// unpack result
	err = n.stakingSMC.Abi.UnpackIntoInterface(&result, "totalSupply", res)
	if err != nil {
		n.lgr.Error("Error unpacking circulating supply error: ", zap.Error(err))
		return nil, err
	}
	return result.TotalSupply, nil
}
