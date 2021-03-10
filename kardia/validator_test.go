// Package kardia
package kardia

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/kardiachain/go-kardia/rpc"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNode_SubscribeStakingEvent(t *testing.T) {
	ctx := context.Background()
	node, err := SetupWSNodeClient()
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

func setupTestNode() (*node, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	node := &node{
		client: rpcClient,
		lgr:    zap.L(),
	}
	if err := node.setupSMC(); err != nil {
		return nil, err
	}
	return node, nil
}

func BenchmarkValidators_List(b *testing.B) {
	n, err := setupTestNode()
	if err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		n.Validators(context.Background())
	}
}

func BenchmarkValidator_CommissionRate(b *testing.B) {
	n, err := setupTestNode()
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
	node, err := SetupNodeClient()
	assert.Nil(t, err)
	hodler := "0x4dAe614b2eA2FaeeDDE7830A2e7fcEDdAE9f9161"
	v, err := node.ValidatorInfo(ctx, hodler)
	assert.Nil(t, err)
	fmt.Println("Validator Details", v)

}

func TestValidator_GetDelegation(t *testing.T) {
	ctx := context.Background()
	node, err := SetupNodeClient()
	assert.Nil(t, err)

	provins := "0xba3D559137Bfe98D3BF0d44D7b34B15A9530c060"
	addr := "0x458892022e66FE0Ef264fE6240EE59fC2FB0A62C"
	stakedAmount, err := node.DelegatorStakedAmount(ctx, provins, addr)
	assert.Nil(t, err)
	fmt.Println("StakedAmount", stakedAmount.String())

}

func TestValidator_List(t *testing.T) {
	ctx := context.Background()
	node, err := SetupNodeClient()
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

		//signInfo, err := node.SigningInfo(ctx, smcAddr.Hex())
		//assert.Nil(t, err)
		//fmt.Println("SignInfo", signInfo)

		//delegatorAddresses, err := node.DelegatorAddresses(ctx, smcAddr.Hex())
		//assert.Nil(t, err)
		//for _, addr := range delegatorAddresses {
		//	if addr.Equal(validator.Signer) {
		//		// Get self delegated info
		//		stakedAmount, err := node.DelegatorStakedAmount(ctx, smcAddr.Hex(), addr.Hex())
		//		assert.Nil(t, err)
		//		fmt.Println("Self-staked", stakedAmount.String())
		//	}
		//}
	}
}

func calculateStats(t *testing.T) {
	ctx := context.Background()
	node, err := SetupNodeClient()
	assert.Nil(t, err)
	totalStaked, err := node.TotalStakedAmount(ctx)
	assert.Nil(t, err)
	fmt.Println("Total staked", totalStaked.String())
}

func Test_GetValidatorsOfDelegator(t *testing.T) {

}
