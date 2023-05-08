package node

import (
	"fmt"

	"math/big"
	"testing"

	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/stretchr/testify/assert"

	"github.com/kardiachain/go-kaiclient/abis"
	"github.com/kardiachain/go-kaiclient/rpc"
)

func TestContractCall(t *testing.T) {
	c, err := rpc.NewClient("https://rpc.kardiachain.io")
	assert.Nil(t, err)
	eth := NewEth(c)
	bookContract, err := eth.NewContract(abis.BookABI, "0xaa1d453082C728D5589Bd4Ca2fF94dC07949B80f")
	assert.Nil(t, err)
	methods := bookContract.AllMethods()
	for i, method := range methods {
		fmt.Printf("Method %d, %s \n", i, method)
	}

	nftSetID := new(big.Int).SetInt64(1)
	//data, err := bookContract.EncodeABI(
	//	"balanceOf",
	//	common.StringToAddress("0xc40295f5a6e4B37CD4FEE8EA1ac8b357559c4799"),
	//	nftSetID,
	//)
	assert.Nil(t, err)
	response, err := bookContract.Call("balanceOf", common.StringToAddress("0xc40295f5a6e4B37CD4FEE8EA1ac8b357559c4799"),
		nftSetID)
	assert.Nil(t, err)
	fmt.Printf("Response: %+v \n", response)
}

func TestCert_BalanceOf(t *testing.T) {
	c, err := rpc.NewClient("https://rpc.kardiachain.io")
	assert.Nil(t, err)
	eth := NewEth(c)
	bookContract, err := eth.NewContract(abis.BookABI, "0xaa1d453082C728D5589Bd4Ca2fF94dC07949B80f")
	assert.Nil(t, err)
	methods := bookContract.AllMethods()
	for i, method := range methods {
		fmt.Printf("Method %d, %s \n", i, method)
	}

	nftSetID := new(big.Int).SetInt64(1)
	//data, err := bookContract.EncodeABI(
	//	"balanceOf",
	//	common.StringToAddress("0xc40295f5a6e4B37CD4FEE8EA1ac8b357559c4799"),
	//	nftSetID,
	//)
	assert.Nil(t, err)
	response, err := bookContract.Call("balanceOf", common.StringToAddress("0xc40295f5a6e4B37CD4FEE8EA1ac8b357559c4799"),
		nftSetID)
	assert.Nil(t, err)
	fmt.Printf("Response: %+v \n", response)
}

func TestCert_MakeTx(t *testing.T) {
	certABI := abis.CertABI
	c, err := rpc.NewClient("https://rpc.kardiachain.io")
	assert.Nil(t, err)
	eth := NewEth(c)
	if err := eth.SetAccount(privateKeyUsedForTest); err != nil {
		t.Fatal(err)
	}

	certContract, err := eth.NewContract(certABI, "0x30020213DCB5a7C7c3120D480D05Ce0aA146Ab91")
	if err != nil {
		t.Fatal(err)
	}

	methods := certContract.AllMethods()
	for i, method := range methods {
		fmt.Printf("Method %d, %s \n", i, method)
	}
	nftSetID := new(big.Int).SetInt64(1)
	data, err := certContract.EncodeABI("mint",
		common.StringToAddress("0xc40295f5a6e4B37CD4FEE8EA1ac8b357559c4799"),
		nftSetID,
	)
	assert.Nil(t, err)
	//call := &types.CallMsg{
	//	From: eth.address,
	//	To:   faucet.Address(),
	//	Data: data,
	//	Gas:  types.NewCallMsgBigInt(big.NewInt(types.MAX_GAS_LIMIT)),
	//}

	gasLimit := uint64(3100000)
	gasPrice := big.NewInt(1000000000)

	fmt.Printf("Estimate gas limit %v\n", gasLimit)
	nonce, err := eth.GetNonce(eth.Address(), nil)
	assert.Nil(t, err)

	txHash, err := eth.SyncSendRawTransaction(
		certContract.Address(),
		big.NewInt(0),
		nonce,
		gasLimit,
		gasPrice,
		data,
	)
	assert.Nil(t, err)
	fmt.Printf("Send faucet tx hash %v\n", txHash)

}
