// Package kardia
package kardia

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/stretchr/testify/assert"
)

var (
	WheelABIJson = `[
	{
		"inputs": [],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "_to",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "_amount",
				"type": "uint256"
			}
		],
		"name": "transferReward",
		"type": "event"
	},
	{
		"inputs": [],
		"name": "TIME",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "day",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_amount",
				"type": "uint256"
			}
		],
		"name": "emergencyWithdrawalKAI",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getBalance",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "maxReward",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "maxSpinDaily",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "pause",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "paused",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256[]",
				"name": "_kaiRange",
				"type": "uint256[]"
			}
		],
		"name": "setKaiRange",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_maxReward",
				"type": "uint256"
			}
		],
		"name": "setMaxReward",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_maxSpinDaily",
				"type": "uint256"
			}
		],
		"name": "setMaxSpinPerDay",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "spin",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "startTime",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "totalReward",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "unpause",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"stateMutability": "payable",
		"type": "receive"
	}
]`
)

func GetWheelAbi() *abi.ABI {
	abiData, _ := abi.JSON(strings.NewReader(WheelABIJson))
	return &abiData
}

func GetSpinRewardFilter() Filter {
	abiData := GetWheelAbi()
	eventCodes := []string{"transferReward"}
	filter, _ := NewFilter(eventCodes, abiData)

	return filter
}

func TestFilter_EventLogs(t *testing.T) {
	rewardFilter := GetSpinRewardFilter()
	ctx := context.Background()
	txHash := "0x10beb023c8f1c9043934064ab7a630a53ff23b423b1cddbc523ab3b818d023ae"
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)
	receipt, err := node.GetTransactionReceipt(ctx, txHash)
	assert.Nil(t, err)

	for _, l := range receipt.Logs {
		l.Data = l.Data[2:]
		events, err := rewardFilter.Events(l)
		if err != nil {
			panic(err.Error())
		}
		for _, event := range events {
			//fmt.Printf("ev: %+v", event)
			amountR := event.Inputs["_amount"]
			fmt.Println("Tx ", txHash, " Found tx amount", amountR)
			//amount = amountR.(*big.Int)
		}
	}
}
