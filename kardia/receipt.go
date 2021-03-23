// Package kardia
package kardia

import (
	"encoding/hex"
	"strings"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
)

type IReceipt interface {
	DecodeInputData(to string, input string) (*FunctionCall, error)
}

//IsKRC20

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

func (n *node) DecodeWithABI(to string, input string, a *abi.ABI) (*FunctionCall, error) {
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
