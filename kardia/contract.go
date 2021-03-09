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
}

func NewContract(abi *abi.ABI, addr common.Address) *Contract {
	c := &Contract{
		Abi:             abi,
		ContractAddress: addr,
	}
	return c
}

type BoundContract struct {
	Abi             *abi.ABI
	ContractAddress common.Address
	Bytecode        string

	*bind.BoundContract
}

func NewBoundContract(node Node, abi *abi.ABI, addr common.Address) *BoundContract {
	c := &BoundContract{
		Abi:             abi,
		ContractAddress: addr,
		BoundContract:   bind.NewBoundContract(addr, *abi, node, node, nil),
	}

	return c
}
