// Package kardia
package kardia

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func Test_node_BlockByHash(t *testing.T) {
	lgr, _ := zap.NewDevelopment()
	ctx := context.Background()
	node, _ := NewNode("https://dev-1.kardiachain.io", lgr)
	b, err := node.BlockByHash(ctx, "0xeb250a3b4efcca94ba29ae9fb5d39b90369bd1d2688e4b8b75fffa11357f3821")
	assert.Nil(t, err)

	fmt.Println("b", b)
}
