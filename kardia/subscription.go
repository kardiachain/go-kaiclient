// Package kardia
package kardia

import (
	"context"

	"github.com/kardiachain/go-kardia/rpc"
	"github.com/kardiachain/go-kardia/types"
)

type ISubscription interface {
	KaiSubscribe(ctx context.Context, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error)

	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (*rpc.ClientSubscription, error)
}

func (n *node) KaiSubscribe(ctx context.Context, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error) {
	return n.client.Subscribe(ctx, "kai", channel, args...)
}

// SubscribeNewHead subscribes to notifications about the current blockchain head
// on the given channel.
func (n *node) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (*rpc.ClientSubscription, error) {
	return n.KaiSubscribe(ctx, ch, "newHeads")
}
