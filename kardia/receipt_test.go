package kardia

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestUnpackLogs(t *testing.T) {
	lgr, _ := zap.NewDevelopment()
	tx := "0x7edbf9c942e0908ce095d975d758894d7a082dc6d6b10281a53a4bb21d718085"
	node, err := NewNode("https://dev.kardiachain.io/", lgr)
	assert.Nil(t, err)
	r, err := node.GetTransactionReceipt(context.Background(), tx)
	assert.Nil(t, err)
	abi, err := KRC721ABI()
	assert.Nil(t, err)
	for _, l := range r.Logs {
		// Process if transfer event
		if l.Topics[0] == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" {

			unpackedLog, err := UnpackLog(l, abi)
			assert.Nil(t, err)
			lgr.Debug("Final logs", zap.Any("UnpackedLog", unpackedLog))
		}

		// Process if mint/burn event

	}
}
