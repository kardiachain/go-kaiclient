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
)

func TestToken_KRC721(t *testing.T) {
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)
	ctx := context.Background()
	r, err := node.GetTransactionReceipt(ctx, "0xf25a2aed31d1eb87dbc151ef0b9f62750b4740416d5c512bca3c9641192d19be")
	assert.Nil(t, err)
	fmt.Println("Receipts", r)
	abi, err := KRC721ABI()
	assert.Nil(t, err)
	for _, l := range r.Logs {
		unpackLog, err := UnpackLog(l, abi)
		if err != nil {
			fmt.Printf("LogsInfo: %+v \n", l)
		} else {
			fmt.Printf("UnpackLogInfo: %+v \n", unpackLog)
		}

	}
}
