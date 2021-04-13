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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewKRC20(t *testing.T) {
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)
	krc20, err := NewKRC20(node, "0x443c6B6240AFDf1D7497aAB64858F5b1363d321B", "0x4f36A53DC32272b97Ae5FF511387E2741D727bdb")
	assert.Nil(t, err)
	fmt.Println("krc20", krc20)
}

func TestHolderBalance(t *testing.T) {
	ctx := context.Background()
	address := "0x4f36A53DC32272b97Ae5FF511387E2741D727bdb"
	pairSMCAddr := "0x0f0524Aa6c70d8B773189C0a6aeF3B01719b0b47"
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)
	krc, err := NewKRC20(node, pairSMCAddr, "")
	assert.Nil(t, err)
	balance, err := krc.HolderBalance(ctx, address)
	assert.Nil(t, err)
	fmt.Println("Balance", balance)
}

func TestA(t *testing.T) {
	logger, err := zap.NewDevelopment()
	node, err := NewNode("wss://ws-dev.kardiachain.io", logger)
	assert.Nil(t, err)
	fmt.Println("Client", node)

}
