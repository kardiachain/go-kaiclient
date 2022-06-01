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

func TestBlock_BlockByHeight(t *testing.T) {
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)
	ctx := context.Background()
	b, err := node.BlockByHeight(ctx, 78785)
	assert.Nil(t, err)
	assert.Equal(t, 78785, b.Height)
}

func TestBlock_BlockByHash(t *testing.T) {
	ctx := context.Background()
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)
	b, err := node.BlockByHash(ctx, "0x5f47a97ac4d8454312430b48f1841af0e839a8fc9f7fc6c00efadb89ea3e2133")
	assert.Nil(t, err)
	fmt.Println("Block", b)
	assert.Equal(t, "0x5f47a97ac4d8454312430b48f1841af0e839a8fc9f7fc6c00efadb89ea3e2133", b.Hash)
}
