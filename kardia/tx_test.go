// Package kardia
package kardia

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/crypto"
	"github.com/kardiachain/go-kardia/types"
	"github.com/stretchr/testify/assert"
)

func TestSendSignedTx(t *testing.T) {
	receivedAddress := common.HexToAddress("0x59173FAF22C3fEd212Ec6B5Ea2E50f7644b614f3")
	pubKey, privKey, err := setupTestAccount()
	assert.Nil(t, err)
	fromAddress := crypto.PubkeyToAddress(*pubKey)

	node, err := setupTestNodeInterface()
	assert.Nil(t, err)
	nonce, err := node.NonceAt(context.Background(), fromAddress.String())
	assert.Nil(t, err)
	balance, err := node.Balance(context.Background(), fromAddress.String())
	assert.Nil(t, err)
	fmt.Println("Balance", balance)
	gasLimit := uint64(3100000)
	gasPrice := big.NewInt(1000000000)
	// Send 1 KAI to from test account to receivedAddress
	var Hydro = big.NewInt(1000000000000000000)
	oneKai := new(big.Int).Mul(new(big.Int).SetInt64(1), Hydro)

	//nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte
	tx := types.NewTransaction(nonce, receivedAddress, oneKai, gasLimit, gasPrice, nil)
	signedTx, err := types.SignTx(types.HomesteadSigner{}, tx, privKey)
	assert.Nil(t, err)

	err = node.SendTransaction(context.Background(), signedTx)
	assert.Nil(t, err)
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}

func TestTx(t *testing.T) {
	n, err := setupTestNodeInstance()
	assert.Nil(t, err)
	tx, err := n.GetTransaction(context.Background(), "0x48bc30dc48af02cf8c39ae744230f4d0725467587c72e5f8f8a53e3113e62e32")
	assert.Nil(t, err)
	fmt.Printf("tx: %+v \n", tx)
}
