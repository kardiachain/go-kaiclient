// Package kardia
package kardia

import (
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
)

type IReceipt interface {
	DecodeInputData(to string, input string) (*FunctionCall, error)
	//DecodeWithABI(to string, input string, a *abi.ABI) (*FunctionCall, error)
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

func UnpackLog(log *Log, smcABI *abi.ABI) (*Log, error) {
	if strings.HasPrefix(log.Data, "0x") {
		log.Data = log.Data[2:]
	}

	event, err := smcABI.EventByID(common.HexToHash(log.Topics[0]))
	if err != nil {
		return nil, err
	}
	argumentsValue := make(map[string]interface{})
	err = unpackLogIntoMap(smcABI, argumentsValue, event.RawName, log)
	if err != nil {
		return nil, err
	}
	//convert address, bytes and string arguments into their hex representations
	for i, arg := range argumentsValue {
		argumentsValue[i] = parseBytesArrayIntoString(arg)
	}
	// append unpacked data
	log.Arguments = argumentsValue
	log.MethodName = event.RawName
	order := int64(1)
	for _, arg := range event.Inputs {
		if arg.Indexed {
			log.ArgumentsName += "index_topic_" + strconv.FormatInt(order, 10) + " "
			order++
		}
		log.ArgumentsName += arg.Type.String() + " " + arg.Name + ", "
	}
	log.ArgumentsName = strings.TrimRight(log.ArgumentsName, ", ")
	return log, nil
}

// UnpackLogIntoMap unpacks a retrieved log into the provided map.
func unpackLogIntoMap(a *abi.ABI, out map[string]interface{}, eventName string, log *Log) error {
	data, err := hex.DecodeString(log.Data)
	if err != nil {
		return err
	}

	if len(data) > 0 {
		if err := a.UnpackIntoMap(out, eventName, data); err != nil {
			return err
		}
	}
	// unpacking indexed arguments
	var indexed abi.Arguments
	for _, arg := range a.Events[eventName].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}

	topicSize := len(log.Topics)
	if topicSize <= 1 {
		return nil
	}
	topics := make([]common.Hash, len(log.Topics)-1)
	for i, topic := range log.Topics[1:] { // exclude the eventID (log.Topic[0])
		topics[i] = common.HexToHash(topic)
	}
	return abi.ParseTopicsIntoMap(out, indexed, topics)
}
