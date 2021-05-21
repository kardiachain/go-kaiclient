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
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/crypto"
	"github.com/kardiachain/go-kardia/rpc"
	"github.com/kardiachain/go-kardia/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSubscription_NewBlockHead(t *testing.T) {
	lgr, err := zap.NewProduction()
	assert.Nil(t, err)
	url := "wss://ws-dev.kardiachain.io/ws"

	node, err := NewNode(url, lgr)
	assert.Nil(t, err)

	headersCh := make(chan *types.Header)
	sub, err := node.SubscribeNewHead(context.Background(), headersCh)
	assert.Nil(t, err, "cannot subscribe")

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headersCh:
			fmt.Println(header.Hash().Hex())
		}
	}
}

func subscribe(n Node, channel interface{}) (*rpc.ClientSubscription, error) {
	args := FilterArgs{Address: []string{"0x9e003D7e05E19514aa76DcFF4EB2443522677288"}}
	sub, err := n.KaiSubscribe(context.Background(), channel, "logs", args)
	if err != nil {
		return nil, err
	}

	return sub, nil

	////rpcClient, err := rpc.Dial("ws://10.10.0.68:8546/ws")
	//rpcClient, err := rpc.Dial("ws://10.10.loo0.251:8550/ws")
	//assert.Nil(t, err, "cannot connect") //NewHeads
	//sub, err := rpcClient.Subscribe(context.Background(), "kai", headersCh, "newHeads")
}

func start(sub *rpc.ClientSubscription, channel chan *FilterLogs) error {
	for {
		select {
		case err := <-sub.Err():
			return err
		case logsEvent := <-channel:
			fmt.Println(logsEvent) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f//
		}
	}
}

func TestSubscription_LogsFilter(t *testing.T) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	lgr, err := zap.NewDevelopment()
	assert.Nil(t, err)
	url := "wss://ws-dev.kardiachain.io/ws"
	//url := "ws://10.10.0.251:8550/ws"
	node, err := NewNode(url, lgr)
	assert.Nil(t, err)

	for {
		lgr.Debug("Start subscribe flow")
		time.Sleep(1 * time.Second)
		logEventCh := make(chan *FilterLogs, 10)
		sub, err := subscribe(node, logEventCh)
		if err != nil {
			// todo: handle close graceful
			lgr.Debug("Cannot subscribe, closed", zap.Error(err))
			return
		}
		err = start(sub, logEventCh)
		switch err {
		case nil:
			continue
		default:
			sub.Unsubscribe()
			lgr.Debug("cannot start", zap.Error(err))
		}
	}

}

func TestSubscription_LogsFilter2(t *testing.T) {
	//startTime := time.Now()
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	lgr, err := zap.NewDevelopment()
	assert.Nil(t, err)
	url := "wss://ws-dev.kardiachain.io/ws"
	//url := "ws://10.10.0.251:8550/ws"
	node, err := NewNode(url, lgr)
	assert.Nil(t, err)
	go func() {
		for {
			_, err := node.LatestBlockNumber(context.Background())
			if err != nil {
				return
			}
			lgr.Debug("Ping by get latest block number", zap.Time("H", time.Now()))
			time.Sleep(10 * time.Second)
		}

	}()

	for {
		lgr.Debug("Start subscribe flow")
		time.Sleep(1 * time.Second)
		logEventCh := make(chan *FilterLogs, 10)
		sub, err := subscribe(node, logEventCh)
		if err != nil {
			// todo: handle close graceful
			lgr.Debug("Cannot subscribe, closed", zap.Error(err))
			return
		}

		err = start(sub, logEventCh)
		switch err {
		case nil:
			continue
		default:
			sub.Unsubscribe()
			lgr.Debug("cannot start", zap.Error(err))
		}
	}
}

func TestSubscription_LogsFilter3(t *testing.T) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	lgr, err := zap.NewDevelopment()
	assert.Nil(t, err)
	url := "wss://ws-dev.kardiachain.io/ws"
	//url := "ws://10.10.0.251:8550/ws"
	node, err := NewNode(url, lgr)
	assert.Nil(t, err)

	startTime := time.Now()
	for {
		lgr.Debug("Start subscribe flow")
		logEventCh := make(chan *FilterLogs, 10)
		sub, err := subscribe(node, logEventCh)
		if err != nil {
			// todo: handle close graceful
			lgr.Debug("Drop here", zap.Duration("Total", time.Now().Sub(startTime)))
			lgr.Debug("Cannot subscribe, closed", zap.Error(err))
			return
		}
		err = start(sub, logEventCh)
		switch err {
		case nil:
			continue
		default:
			sub.Unsubscribe()
			lgr.Debug("cannot start", zap.Error(err))
		}
	}

}

const testAbi = `[
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "test2",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "uint256",
				"name": "index",
				"type": "uint256"
			}
		],
		"name": "EventOrderTest",
		"type": "event"
	},
	{
		"constant": false,
		"inputs": [
			{
				"internalType": "uint256",
				"name": "n",
				"type": "uint256"
			}
		],
		"name": "StressTestEvents",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "test2Address",
		"outputs": [
			{
				"internalType": "contract Test2",
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	}
]`

func pumpSubscribe(n Node, channel interface{}) (*rpc.ClientSubscription, error) {
	args := FilterArgs{Address: []string{"0x11A8CC3A9e8e8AF4667C26B7CC68873C86Ef43B0", "0x60e3a154547E3D940EA7af651A65E838d7672Ba2"}} //test1, test2
	//{"jsonrpc":"2.0", "id": 1, "method": "kai_subscribe", "params": ["logs", {"address": ["0x11A8CC3A9e8e8AF4667C26B7CC68873C86Ef43B0","0x60e3a154547E3D940EA7af651A65E838d7672Ba2"]}]}
	sub, err := n.KaiSubscribe(context.Background(), channel, "logs", args)
	if err != nil {
		return nil, err
	}

	return sub, nil

	////rpcClient, err := rpc.Dial("ws://10.10.0.68:8546/ws")
	//rpcClient, err := rpc.Dial("ws://10.10.loo0.251:8550/ws")
	//assert.Nil(t, err, "cannot connect") //NewHeads
	//sub, err := rpcClient.Subscribe(context.Background(), "kai", headersCh, "newHeads")
}

var (
	test1SmcAddr        = "0x11A8CC3A9e8e8AF4667C26B7CC68873C86Ef43B0"
	gasLimit            = new(big.Int).SetInt64(39999999)
	gasPrice            = new(big.Int).SetInt64(1000000000)
	nonce        uint64 = 3533
	total               = 0
)

func TestSubscription_LogsFilter_Pump(t *testing.T) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	lgr, err := zap.NewDevelopment()
	assert.Nil(t, err)

	httpNode, err := setupTestNodeInterface()
	assert.Nil(t, err)
	r := strings.NewReader(testAbi)
	abiData, err := abi.JSON(r)
	assert.Nil(t, err)
	smc := Contract{
		Abi:             &abiData,
		ContractAddress: common.HexToAddress("0x11A8CC3A9e8e8AF4667C26B7CC68873C86Ef43B0"),
	}

	url := "wss://ws-dev.kardiachain.io"
	//url := "ws://10.10.0.251:8550/ws"
	node, err := NewNode(url, lgr)
	assert.Nil(t, err)

	lgr.Debug("Start subscribe flow")
	time.Sleep(1 * time.Second)
	logEventCh := make(chan *FilterLogs, 10)
	sub, err := pumpSubscribe(node, logEventCh)
	if err != nil {
		// todo: handle close graceful
		lgr.Debug("Cannot subscribe, closed", zap.Error(err))
		return
	}
	var (
		previousAddress       = "0x60e3a154547E3D940EA7af651A65E838d7672Ba2"
		previous        int64 = 0
	)
	//pubKey, _, err := setupTestAccount()
	//assert.Nil(t, err)
	//fromAddress := crypto.PubkeyToAddress(*pubKey)
	//nonce, err = httpNode.NonceAt(context.Background(), fromAddress.String())
	//fmt.Printf("@@@@@@@@@@@@@@@@@@ NonceAt error %v\n", nonce)
	//if err != nil {
	//	fmt.Printf("@@@@@@@@@@@@@@@@@@ NonceAt error %v\n", err)
	//}
	recursiveCall(httpNode, smc)
	for {
		select {
		case err := <-sub.Err():
			lgr.Debug("subscribe err", zap.Error(err))
		case log := <-logEventCh:
			total++
			if (log.LogIndex-1 == previous || log.LogIndex == 1) && !strings.EqualFold(log.Address, previousAddress) {
				if log.LogIndex == 1 {
					t.Log("new block", "previous", previous, "total", total, "blockHeight", log.BlockHeight)
				} else if total == 4000 {
					t.Log("total", total, "blockHeight", log.BlockHeight)
				}
				previous = log.LogIndex
				previousAddress = log.Address
			} else {
				t.Fatal("wrong order of logs", "logIndex", log.LogIndex, "previous", previous)
			}
			//t.Log("event details", "address", log.Address, "index", log.LogIndex, "total", total, "hash", log.TransactionHash)
		}
	}
}

func recursiveCall(node Node, smc Contract) {
	payload, err := smc.Abi.Pack("StressTestEvents", new(big.Int).SetInt64(20))
	if err != nil {
		fmt.Printf("@@@@@@@@@@@@@@@@@@ Pack error %v\n", err)
	}

	_, privKey, err := setupSender()
	if err != nil {
		fmt.Printf("@@@@@@@@@@@@@@@@@@ setupSender error %v\n", err)
	}

	nonce, err = node.NonceAt(context.Background(), "0xFBD5e2aFB7C0a7862b06964e29E676bf02183256")
	if err != nil {
		fmt.Println("Err", err)
		return
	}
	fmt.Println("Nonce", nonce)
	for i := 1; i <= 200; i++ {
		tx := types.NewTransaction(nonce, common.HexToAddress(test1SmcAddr), new(big.Int).SetInt64(0), gasLimit.Uint64(), gasPrice, payload)
		nonce++
		signedTx, err := types.SignTx(types.HomesteadSigner{}, tx, privKey)
		if err != nil {
			fmt.Printf("@@@@@@@@@@@@@@@@@@ SignTx error %v\n", err)
		}

		err = node.SendTransaction(context.Background(), signedTx)
		if err != nil {
			fmt.Printf("@@@@@@@@@@@@@@@@@@ SendTransaction error %v\n", err)
		}
	}
	fmt.Println("new block", "total", total)
}

func setupSender() (*ecdsa.PublicKey, *ecdsa.PrivateKey, error) {
	privateKey, err := crypto.HexToECDSA("8063515889cb660ab1c25f48470b883c75515a83374797aef99f7f81c18ecd11")
	if err != nil {
		return nil, nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, err
	}

	return publicKeyECDSA, privateKey, nil
}
