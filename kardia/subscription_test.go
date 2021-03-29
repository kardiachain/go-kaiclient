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

	args := FilterArgs{Address: []string{"0xe6D4dB026810ad0b405a1e48E8DfFF2509Bc6b0A", "0x50a26DF56fC91eECF7f25D52eFB4eFAB56Dacf08", "0x16B10970D6712D1D87aA18B9770b0Abd969dC079"}}
	logEventCh := make(chan *FilterLogs)
	_, err = node.KaiSubscribe(context.Background(), logEventCh, "logs", args)
	assert.Nil(t, err, "cannot subscribe")

	////rpcClient, err := rpc.Dial("ws://10.10.0.68:8546/ws")
	//rpcClient, err := rpc.Dial("ws://10.10.loo0.251:8550/ws")
	//assert.Nil(t, err, "cannot connect") //NewHeads

	for {
		select {
		case logsEvent := <-logEventCh:
			fmt.Println(logsEvent) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
		}
	}
	//sub, err := node.SubscribeNewHead(context.Background(), headers)
	//assert.Nil(t, err)
	//
}
