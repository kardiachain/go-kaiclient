package node

import (
	"fmt"
	"testing"

	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/stretchr/testify/assert"

	"github.com/kardiachain/go-kaiclient/rpc"
)

var privateKeyUsedForTest = "39da0604ce1db485db2e9e346dfcebee30681b649c6e925de2a9a9cba71c8f13"

func TestSignText(t *testing.T) {
	eth := NewEth(nil)
	if err := eth.SetAccount(privateKeyUsedForTest); err != nil {
		t.Fatal(err)
	}

	signature, err := eth.SignText([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("signature %x\n", signature)
}

func TestReceipt(t *testing.T) {
	c, err := rpc.NewClient("https://dev.kardiachain.io")
	if err != nil {
		t.Fatal(err)
	}
	eth := NewEth(c)
	hash := common.HexToHash("0x326b0a98c45182d3df381c87055c6a96c942b20bdb5a7fd7d64935112fd789c1")
	r, err := eth.GetTransactionReceipt(hash)
	assert.Nil(t, err)
	fmt.Println("r", r)
}
