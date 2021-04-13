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
