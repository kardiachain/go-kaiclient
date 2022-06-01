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
	"fmt"
	"log"
	"testing"
	"time"

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
	args := FilterArgs{Address: []string{"0x42d3400560F66A15F6D1345b894A854E5277270a"}}
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
	url := "ws://10.10.0.251:8550/ws"
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
	url := "wss://ws.kardiachain.io/ws"
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
