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
	node, err := SetupNodeClient()
	assert.Nil(t, err)

	stakedAmount, err := node.TotalStakedAmount(ctx)
	assert.Nil(t, err)
	fmt.Println("Staked amount", stakedAmount)

}

func TestStaking_ValidatorSMCAddresses(t *testing.T) {
	ctx := context.Background()
	node, err := SetupNodeClient()
	assert.Nil(t, err)
	validatorSMCAddresses, err := node.ValidatorSMCAddresses(ctx)
	assert.Nil(t, err)
	fmt.Println("List", validatorSMCAddresses)

}
