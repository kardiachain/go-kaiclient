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
	if reflect.TypeOf(v).Kind() == reflect.Array {
		arr := v.([32]byte)
		slice := arr[:]
		// convert any array of uint8 into a hex string
		if reflect.TypeOf(slice).Elem().Kind() == reflect.Uint8 {
			return common.Bytes(slice).String()
		} else {
			// otherwise recursively check other arguments
			return parseBytesArrayIntoString(v)
		}
	} else if reflect.TypeOf(v).Kind() == reflect.Ptr {
		// convert big.Int to string to avoid overflowing
		if value, ok := v.(*big.Int); ok {
			return value.String()
		}
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
