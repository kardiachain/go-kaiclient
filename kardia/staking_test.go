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

func TestNode_TotalStakedAmount(t *testing.T) {
	ctx := context.Background()
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)

	stakedAmount, err := node.TotalStakedAmount(ctx)
	assert.Nil(t, err)
	fmt.Println("Staked amount", stakedAmount)

}

func TestStaking_ValidatorSMCAddresses(t *testing.T) {
	ctx := context.Background()
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)
	validatorSMCAddresses, err := node.ValidatorSMCAddresses(ctx)
	assert.Nil(t, err)
	fmt.Println("List", validatorSMCAddresses)
}
