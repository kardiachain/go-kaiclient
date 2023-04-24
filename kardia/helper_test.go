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
	"fmt"
	"math/big"
	"testing"

	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/crypto"
	"github.com/stretchr/testify/assert"
)

func TestGenerateNewWallet(t *testing.T) {
	address, privKey, err := GenerateWallet()
	assert.Nil(t, err)
	privateKeyBytes := crypto.FromECDSA(&privKey)
	privateKeyStr := common.Encode(privateKeyBytes)
	fmt.Println("WalletAddress", address.Hex())
	fmt.Println("WalletPrivateKey", privateKeyStr)
}

func TestFloatToBigInt(t *testing.T) {
	amount := FloatToBigInt(0.2, 18)
	expectedAmount, _ := new(big.Int).SetString("200000000000000000", 10)
	assert.Equal(t, expectedAmount, amount)

}
