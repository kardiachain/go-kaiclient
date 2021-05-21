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

	"github.com/kardiachain/go-kardia/lib/crypto"
	"github.com/stretchr/testify/assert"
)

func TestBoundContract_Deploy(t *testing.T) {
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)
	//r := strings.NewReader(smc.KRC20ABI)
	//abiData, err := abi.JSON(r)
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
	gasPrice := big.NewInt(1000000000)
	auth := NewKeyedTransactor(privateKey)
	auth.Nonce = nonce
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice

	smcAddress, txHash, err := node.DeployKRC20(auth)
	assert.Nil(t, err)
	fmt.Println("SMC Addr", smcAddress.String())
	fmt.Println("TxHash", txHash.String())
}
