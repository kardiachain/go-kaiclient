// Package kardia
package kardia

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/rpc"
	"github.com/kardiachain/go-kardia/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSubscription_NewBlockHead(t *testing.T) {
	lgr, err := zap.NewDevelopment()
	assert.Nil(t, err)
	url := "wss://ws-dev.kardiachain.io/ws"

	node, err := NewNode(url, lgr)
	assert.Nil(t, err)

	headersCh := make(chan *types.Header)
	sub, err := node.SubscribeNewHead(context.Background(), headersCh)
	assert.Nil(t, err, "cannot subscribe")

	////rpcClient, err := rpc.Dial("ws://10.10.0.68:8546/ws")
	//rpcClient, err := rpc.Dial("ws://10.10.loo0.251:8550/ws")
	//assert.Nil(t, err, "cannot connect") //NewHeads
	//sub, err := rpcClient.Subscribe(context.Background(), "kai", headersCh, "newHeads")

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headersCh:
			fmt.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
		}
	}
	//sub, err := node.SubscribeNewHead(context.Background(), headers)
	//assert.Nil(t, err)
	//
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

func TestSubscription_LogsFilterEvent(t *testing.T) {
	lgr, err := zap.NewDevelopment()
	assert.Nil(t, err)
	//url := "wss://ws-dev.kardiachain.io/ws"
	url := "ws://10.10.0.251:8550/ws"
	node, err := NewNode(url, lgr)
	assert.Nil(t, err)

	args := FilterArgs{Address: []string{"0x0f0524Aa6c70d8B773189C0a6aeF3B01719b0b47", "0xA37888A6FF1B0A5347D8355DE05A044F2981a959", "0x0059DC83f10E32eB653C77942C821ca1Be7D89bB", "0xF742b790f36A4FB2DCbE1795573a0F574fc69F32", "0x2d5407B96950B99d3b67a2364DE6590F02BEBde5", "0xeA80f362fe37616b5922E40A5F6f912045D5e47c", "0xB643f79b5b9eb20A8a7CC93721b83B041E1c6048", "0x034F0B3C54bF9430344E2afC80C71F6A79f11a9B", "0x0FeBeC85D9279097c893BfF7929C95aF8A334121", "0x4C4980B121810B6fa4bc7cAE195609414a4fd213"}}
	logEventCh := make(chan *Log)
	_, err = node.KaiSubscribe(context.Background(), logEventCh, "logs", args)
	assert.Nil(t, err, "cannot subscribe")

	////rpcClient, err := rpc.Dial("ws://10.10.0.68:8546/ws")
	//rpcClient, err := rpc.Dial("ws://10.10.loo0.251:8550/ws")
	//assert.Nil(t, err, "cannot connect") //NewHeads
	r := bytes.NewReader([]byte(PairABI))
	pairSMCABI, err := abi.JSON(r)
	if err != nil {
		return
	}
	for {
		select {
		case logEvent := <-logEventCh:
			fmt.Printf("Events: %+v \n", logEvent)
			logEvent.Data = logEvent.Data[2:]
			lDetails, err := UnpackLog(logEvent, &pairSMCABI)
			if err != nil {
				fmt.Println("Cannot unpack ", err.Error())
				return
			}
			fmt.Printf("Log Details: %+v \n", lDetails)
		}
	}
}
