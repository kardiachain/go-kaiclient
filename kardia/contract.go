// Package kardia
package kardia

import (
	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/abi/bind"
	"github.com/kardiachain/go-kardia/lib/common"
)

type Contract struct {
	Abi             *abi.ABI
	ContractAddress common.Address
	Bytecode        string

	*bind.BoundContract
}

func NewContract(node Node, abi *abi.ABI, addr common.Address) *Contract {
	bc := bind.NewBoundContract(addr, *abi, node, node, nil)
	c := &Contract{
		Abi:             abi,
		ContractAddress: addr,
		BoundContract:   bc,
	}

	return c
}
