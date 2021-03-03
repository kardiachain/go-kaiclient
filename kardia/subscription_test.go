// Package kardia
package kardia

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/kardiachain/go-kardia/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSubscription_NewBlockHead(t *testing.T) {
	lgr, err := zap.NewDevelopment()
	assert.Nil(t, err)
	url := "ws://10.10.0.68:8546/ws"
	//url := "ws://10.10.0.251:8550/ws"
	node, err := NewNode(url, lgr)
	assert.Nil(t, err)

	headersCh := make(chan *types.Header)
	sub, err := node.SubscribeNewHead(context.Background(), headersCh)
	assert.Nil(t, err, "cannot subscribe")

	////rpcClient, err := rpc.Dial("ws://10.10.0.68:8546/ws")
	//rpcClient, err := rpc.Dial("ws://10.10.0.251:8550/ws")
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
