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
	"encoding/hex"
	"go.uber.org/zap"
	"strconv"
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

func UnpackLog(log *Log, smcABI *abi.ABI) (*Log, error) {
	if strings.HasPrefix(log.Data, "0x") {
		log.Data = log.Data[2:]
	}
	log.Address = CheckAddress(log.Address)
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
	lgr, _ := zap.NewDevelopment()
	data, err := hex.DecodeString(log.Data)
	if err != nil {
		return err
	}

	if len(data) > 0 {
		if err := a.UnpackIntoMap(out, eventName, data); err != nil {
			return err
		}
	}
	lgr.Info("Event Name", zap.String("Event", eventName), zap.Any("Inputs", a.Events[eventName].Inputs))
	// unpacking indexed arguments
	var indexed abi.Arguments
	for _, arg := range a.Events[eventName].Inputs {
		lgr.Info("Args", zap.Any("Arg", arg))
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	lgr.Info("Indexed", zap.Any("IndexedSize", len(indexed)), zap.Any("Indexed", indexed))

	topicSize := len(log.Topics)
	if topicSize <= 1 {
		return nil
	}
	topics := make([]common.Hash, len(log.Topics)-1)
	for i, topic := range log.Topics[1:] { // exclude the eventID (log.Topic[0])
		topics[i] = common.HexToHash(topic)
	}
	lgr.Info("Topics", zap.Any("TopicSize", len(topics)), zap.Any("Topics", topics))
	return abi.ParseTopicsIntoMap(out, indexed, topics)
}
