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

	"go.uber.org/zap"
)

type IDelegator interface {
	UnbondedRecords(ctx context.Context, validatorSMCAddress, delegatorAddress string) (*UnbondedRecord, error)
}

func (n *node) UnbondedRecords(ctx context.Context, validatorSMCAddress, delegatorAddress string) (*UnbondedRecord, error) {
	lgr := n.lgr.With(zap.String("method", "UnbondedRecords"))
	payload, err := n.validatorSMC.Abi.Pack("getUBDEntries", delegatorAddress)
	if err != nil {
		lgr.Error("Error packing UDB entry payload: ", zap.Error(err))
		return nil, err
	}
	res, err := n.KardiaCall(ctx, ConstructCallArgs(validatorSMCAddress, payload))
	if err != nil {
		lgr.Error("GetUDBEntry KardiaCall error: ", zap.Error(err))
		return nil, err
	}
	if len(res) == 0 {
		return nil, ErrEmptyList
	}

	var result *UnbondedRecord
	// unpack result
	err = n.validatorSMC.Abi.UnpackIntoInterface(&result, "getUBDEntries", res)
	if err != nil {
		lgr.Error("Error unpacking UDB entry: ", zap.Error(err))
		return nil, err
	}

	return result, nil
}
