// Package kardia
package kardia

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNode_SubscribeStakingEvent(t *testing.T) {
	ctx := context.Background()
	node, err := setupWSNodeClient()
	assert.Nil(t, err)
	ch := make(chan interface{})

	type FilterArgs struct {
		From    uint64
		To      uint64
		Address []string
		Topics  []string
	}
	args := []interface{}{
		"logs",
		FilterArgs{Address: []string{"0x42d3400560F66A15F6D1345b894A854E5277270a"}},
	}
	sub, err := node.KaiSubscribe(ctx, ch, "subscribe", args)
	assert.Nil(t, err)
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case event := <-ch:
			fmt.Println("Event", event) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
		}
	}
}

func BenchmarkValidators_List(b *testing.B) {
	n, err := setupTestNodeInstance()
	if err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		n.Validators(context.Background())
	}
}

func BenchmarkValidator_CommissionRate(b *testing.B) {
	n, err := setupTestNodeInstance()
	if err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		now := time.Now()
		if _, err := n.ValidatorCommission(context.Background(), "0x4dAe614b2eA2FaeeDDE7830A2e7fcEDdAE9f9161"); err != nil {
			return
		}
		fmt.Println("Total", time.Now().Sub(now))
	}
}

func TestValidator_Details(t *testing.T) {
	ctx := context.Background()
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)
	hodler := "0xdC4A94805f449A64B27B589233C49d87eE99fBBc"
	v, err := node.ValidatorInfo(ctx, hodler)
	assert.Nil(t, err)
	fmt.Printf("Validator Details: %+v\n", v)
	fmt.Printf("Validator Details: %+v\n", v.SigningInfo)

}

func TestValidator_GetDelegation(t *testing.T) {
	ctx := context.Background()
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)

	vSMC := "0xdC4A94805f449A64B27B589233C49d87eE99fBBc"
	addr := "0x458892022e66FE0Ef264fE6240EE59fC2FB0A62C"
	stakedAmount, err := node.DelegatorStakedAmount(ctx, vSMC, addr)
	assert.Nil(t, err)
	fmt.Println("StakedAmount", stakedAmount.String())

	rewardAMount, err := node.DelegationRewards(ctx, vSMC, addr)
	assert.Nil(t, err)
	fmt.Println("RewardAmount", rewardAMount.String())
}

func TestValidator_List(t *testing.T) {
	ctx := context.Background()
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)

	totalStakedAmount, err := node.TotalStakedAmount(ctx)
	assert.Nil(t, err)
	fmt.Println("TotalStakedAmount", totalStakedAmount.String())
	validatorSMCAddresses, err := node.ValidatorSMCAddresses(ctx)
	assert.Nil(t, err)
	for _, smcAddr := range validatorSMCAddresses {
		// Get basic info
		nValidator, err := node.ValidatorInfo(ctx, smcAddr.Hex())
		assert.Nil(t, err)
		fmt.Printf("ValidatorInfo: %+v \n", nValidator)
		commission, err := node.ValidatorCommission(ctx, smcAddr.Hex())
		assert.Nil(t, err)
		fmt.Println("Commission", commission)
		selfStakedAmount, err := node.DelegatorStakedAmount(ctx, smcAddr.Hex(), nValidator.Signer.Hex())
		assert.Nil(t, err)
		fmt.Println("Self staked amount", selfStakedAmount)
		delegatorAddresses, err := node.DelegatorAddresses(ctx, smcAddr.Hex())
		assert.Nil(t, err)
		fmt.Println("Delegators size", len(delegatorAddresses))

	}
}

func Test_GetValidatorsOfDelegator(t *testing.T) {

}

func Test_ValidatorSets(t *testing.T) {
	ctx := context.Background()
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)

	sets, err := node.ValidatorSets(ctx)
	fmt.Println("Set length", len(sets))
	assert.Nil(t, err)
	for _, add := range sets {
		fmt.Println("Address", add.String())
	}
}

func Test_SMCAddressOfValidator(t *testing.T) {
	ctx := context.Background()
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)

	sets, err := node.SMCAddressOfValidator(ctx, "0x5CdF7E0bBF0C53b5f4e612Fa66f0E60169e3a006")
	assert.Nil(t, err)
	fmt.Println("Set", sets)
} //

func Test_ValidatorAddressOfSMC(t *testing.T) {
	ctx := context.Background()
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)

	sets, err := node.ValidatorAddressOfSMC(ctx, "0x50a26DF56fC91eECF7f25D52eFB4eFAB56Dacf08")
	assert.Nil(t, err)
	fmt.Println("Set", sets)
}

func TestValidator_DelegatorsWithShare(t *testing.T) {
	ctx := context.Background()
	node, err := setupTestNodeInterface()
	assert.Nil(t, err)

	sets, err := node.DelegatorsWithShare(ctx, "0x7D776FEcB1Ecb439aEb1324096f8bd1e6BCd5C03")
	assert.Nil(t, err)
	for _, s := range sets {
		fmt.Printf("Addresss: %s | Share: %s \n", s.Address.String(), s.Share.String())
	}
}
