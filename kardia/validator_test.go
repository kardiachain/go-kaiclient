// Package kardia
package kardia

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode_Validators(t *testing.T) {
	ctx := context.Background()
	node, err := SetupNodeClient()
	assert.Nil(t, err)
	validators, err := node.Validators(ctx)
	assert.Nil(t, err)
	for _, v := range validators {
		fmt.Println("Address", v.Signer.Hex())
		fmt.Println("SMCAddress", v.SMCAddress.Hex())
		fmt.Printf("V Detail: %+v\n ", v)
		fmt.Printf("Sign: %+v \n", v.SigningInfo)
		for _, d := range v.Delegators {
			fmt.Printf("Delegator: %+v \n", d)
		}

	}
}

func TestNode_Validator(t *testing.T) {
	ctx := context.Background()
	node, err := SetupNodeClient()
	assert.Nil(t, err)
	//address := "0xFBD5e2aFB7C0a7862b06964e29E676bf02183256"
	address := "0xf151515fa44527E203Cb457086cDa630da80c4b8" //SMC
	validator, err := node.Validator(ctx, address)
	assert.Nil(t, err)
	fmt.Printf("validator: %+v \n", validator)

}
