// Package kardia
package kardia

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
)

type IContract interface {
	StakingContact(ctx context.Context) *Contract
	ValidatorContact(ctx context.Context) *Contract

	//DecodeLog(ctx context.Context, smcABI *abi.ABI, log *Log) error
	//EstimateGas(ctx context.Context) (uint64, error)
}

// DecodeInputData returns decoded transaction input data if it match any function in staking and validator contract.
func (n *node) DecodeInputData(to string, input string) (*FunctionCall, error) {
	// return nil if input data is too short
	if len(input) <= 2 {
		return nil, nil
	}
	data, err := hex.DecodeString(strings.TrimLeft(input, "0x"))
	if err != nil {
		return nil, err
	}
	sig := data[0:4] // get the function signature (first 4 bytes of input data)
	var (
		a      *abi.ABI
		method *abi.Method
	)
	// check if the to address is staking contract, then we search for staking method in staking contract ABI
	if n.stakingSMC.ContractAddress.Equal(common.HexToAddress(to)) {
		a = n.stakingSMC.Abi
		method, err = n.stakingSMC.Abi.MethodById(sig)
		if err != nil {
			return nil, err
		}
	} else { // otherwise, search for a validator method
		a = n.validatorSMC.Abi
		method, err = n.validatorSMC.Abi.MethodById(sig)
		if err != nil {
			return nil, err
		}
	}
	// exclude the function signature, only decode and unpack the arguments
	var body []byte
	if len(data) <= 4 {
		body = []byte{}
	} else {
		body = data[4:]
	}
	args, err := n.getInputArguments(a, method.Name, body)
	if err != nil {
		return nil, err
	}
	arguments := make(map[string]interface{})
	err = args.UnpackIntoMap(arguments, body)
	if err != nil {
		return nil, err
	}
	// convert address, bytes and string arguments into their hex representations
	for i, arg := range arguments {
		arguments[i] = parseBytesArrayIntoString(arg)
	}
	return &FunctionCall{
		Function:   method.String(),
		MethodID:   "0x" + hex.EncodeToString(sig),
		MethodName: method.Name,
		Arguments:  arguments,
	}, nil
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

//func (n *node) EstimateGas(ctx context.Context) (uint64, error) {
//	var result uint64
//	args := SMCCallArgs{
//		From:     "0x",
//		Gas:      3000000,
//		GasPrice: big.NewInt(1),
//	}
//	err := n.client.CallContext(ctx, &result, "kai_estimateGas", args)
//	return result, err
//}
