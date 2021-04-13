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
	"fmt"
	"math/big"
	"reflect"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
)

type IContract interface {
	StakingContact(ctx context.Context) *Contract
	ValidatorContact(ctx context.Context) *Contract

	//DecodeLog(ctx context.Context, smcABI *abi.ABI, log *Log) error
	//EstimateGas(ctx context.Context) (uint64, error)
}

// parseBytesArrayIntoString is a utility function. It converts address, bytes and string arguments into their hex representation.
func parseBytesArrayIntoString(v interface{}) interface{} {
	vType := reflect.TypeOf(v).Kind()
	switch vType {
	case reflect.Array:
		//common.Address{}
		addr, ok := v.(common.Address)
		if ok {
			return common.Bytes(addr[:]).String()
		}

		hash, ok := v.([32]byte)
		if ok {
			return common.Bytes(hash[:]).String()
		}
		return v
	case reflect.Ptr:
		if value, ok := v.(*big.Int); ok {
			return value.String()
		}
	default:
		return v
	}
	return v
}

// getInputArguments get input arguments of a contract call
func (n *node) getInputArguments(a *abi.ABI, name string, data []byte) (abi.Arguments, error) {
	var args abi.Arguments
	if method, ok := a.Methods[name]; ok {
		if len(data)%32 != 0 {
			return nil, fmt.Errorf("abi: improperly formatted output: %s - Bytes: [%+v]", string(data), data)
		}
		args = method.Inputs
	}
	if args == nil {
		return nil, ErrMethodNotFound
	}
	return args, nil
}

func (n *node) StakingContact(ctx context.Context) *Contract {
	return n.stakingSMC
}

func (n *node) ValidatorContact(ctx context.Context) *Contract {
	return n.validatorSMC
}

type Contract struct {
	Abi             *abi.ABI
	Bytecode        string
	ContractAddress common.Address
	OwnerAddress    common.Address
}

func (c *Contract) SetBytecode(bytecode string) {
	c.Bytecode = bytecode
}

func (c *Contract) SetOwnerAddress(address string) {
	c.OwnerAddress = common.HexToAddress(address)
}

func NewContract(abi *abi.ABI, addr common.Address) *Contract {
	c := &Contract{
		Abi:             abi,
		ContractAddress: addr,
	}
	return c
}
