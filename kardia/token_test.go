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
	"math/big"
	"testing"

	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/crypto"
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
	pairSMCAddr := "0x42020b000d2192082e058ff917d79d38d69D0E7e"
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)
	krc, err := NewKRC20(node, pairSMCAddr, "")
	assert.Nil(t, err)
	balance, err := krc.HolderBalance(ctx, address)
	assert.Nil(t, err)
	fmt.Println("Balance", balance)
}

func TestKRC20_Transfer(t *testing.T) {
	tokenSMCAddr := "0xEC8182CeC08E1588237cDb4bd4E049B8d09Ae821"
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)
	krc, err := NewKRC20(node, tokenSMCAddr, "")
	assert.Nil(t, err)
	pubKey, privateKey, err := setupTestAccount()
	assert.Nil(t, err)
	fromAddress := crypto.PubkeyToAddress(*pubKey)
	nonce, err := node.NonceAt(context.Background(), fromAddress.String())
	assert.Nil(t, err)
	balance, err := node.Balance(context.Background(), fromAddress.String())
	assert.Nil(t, err)
	fmt.Println("Balance", balance)
	gasLimit := uint64(3100000)
	gasPrice := big.NewInt(10000000000)
	auth := NewKeyedTransactor(privateKey)
	auth.Nonce = nonce
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice

	to := common.HexToAddress("0x700460B9972515eC0C7A6709a09C69c30867A704")
	amount, _ := new(big.Int).SetString("1000000000000000000", 10)
	txHash, err := krc.Transfer(context.Background(), auth, to, amount)
	assert.Nil(t, err)
	fmt.Println("txHash", txHash)
}

func TestA(t *testing.T) {
	logger, err := zap.NewDevelopment()
	node, err := NewNode("wss://ws-dev.kardiachain.io", logger)
	assert.Nil(t, err)
	fmt.Println("Client", node)

}
