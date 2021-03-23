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

func (c *BoundContract) DeployKRC20(auth *bind.TransactOpts) (common.Address, common.Hash, error) {
	parsed, err := abi.JSON(strings.NewReader(smc.KRC20ABI))
	if err != nil {
		return common.Address{}, common.Hash{}, err
	}
	address, tx, _, err := bind.DeployContract(auth, parsed, common.FromHex(smc.KRC20Bytecode), c.node)
	if err != nil {
		return common.Address{}, common.Hash{}, err
	}
	return address, tx.Hash(), nil
}
