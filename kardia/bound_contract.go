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
	"strings"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/abi/bind"
	"github.com/kardiachain/go-kardia/lib/common"

	"github.com/kardiachain/go-kaiclient/kardia/smc"
)

type BoundContract struct {
	Abi             *abi.ABI
	ContractAddress common.Address
	Bytecode        string
	node            Node

	*bind.BoundContract
}

func NewBoundContract(node Node, abi *abi.ABI, addr common.Address) *BoundContract {
	c := &BoundContract{
		Abi:             abi,
		ContractAddress: addr,
		node:            node,
		BoundContract:   bind.NewBoundContract(addr, *abi, node, node, nil),
	}

	return c
}

func (n *node) DeployKRC20(auth *bind.TransactOpts) (common.Address, common.Hash, error) {
	parsed, err := abi.JSON(strings.NewReader(smc.KRC20ABI))
	if err != nil {
		return common.Address{}, common.Hash{}, err
	}
	address, tx, _, err := bind.DeployContract(auth, parsed, common.FromHex(smc.KRC20Bytecode), n)
	if err != nil {
		return common.Address{}, common.Hash{}, err
	}
	return address, tx.Hash(), nil
}
