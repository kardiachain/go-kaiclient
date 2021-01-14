// Package kardia
package kardia

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/zap"
)

func Test_node_BlockByHeight(t *testing.T) {
	lgr, _ := zap.NewDevelopment()
	ctx := context.Background()
	node, _ := NewNode("https://dev-1.kardiachain.io", lgr)
	b, err := node.BlockByHeight(ctx, 78785)
	if err != nil {
		return
	}

	fmt.Println("b", b)
}
