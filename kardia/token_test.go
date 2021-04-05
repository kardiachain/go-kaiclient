// Package kardia
package kardia

import (
	"context"
	"fmt"
	"testing"

	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/stretchr/testify/assert"
)

func TestNewKRC20(t *testing.T) {
	node, err := SetupNodeClient()
	assert.Nil(t, err)
	krc20, err := NewKRC20(node, "0x443c6B6240AFDf1D7497aAB64858F5b1363d321B", "0x4f36A53DC32272b97Ae5FF511387E2741D727bdb")
	assert.Nil(t, err)
	fmt.Println("krc20", krc20)
}

func TestHolderBalance(t *testing.T) {
	ctx := context.Background()
	address := "0x4f36A53DC32272b97Ae5FF511387E2741D727bdb"
	pairSMCAddr := "0x0f0524Aa6c70d8B773189C0a6aeF3B01719b0b47"
	node, err := SetupNodeClient()
	assert.Nil(t, err)
	krc, err := NewKRC20(node, pairSMCAddr, "")
	assert.Nil(t, err)
	balance, err := krc.HolderBalance(ctx, common.HexToAddress(address))
	assert.Nil(t, err)
	fmt.Println("Balance", balance)
}
