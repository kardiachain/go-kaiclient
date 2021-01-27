// Package kardia
package kardia

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	wheelABI = `[
	{
		"inputs": [],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"inputs": [],
		"name": "TIME",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_amount",
				"type": "uint256"
			}
		],
		"name": "emergencyWithdrawalKAI",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getBalance",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "maxReward",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"name": "numberOfSpin",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"name": "reward",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256[]",
				"name": "_kaiRange",
				"type": "uint256[]"
			}
		],
		"name": "setKaiRange",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_maxReward",
				"type": "uint256"
			}
		],
		"name": "setMaxReward",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_addr",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "_number",
				"type": "uint256"
			}
		],
		"name": "setNumberOfSpin",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "spin",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "startTime",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "totalReward",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "withdrawReward",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"stateMutability": "payable",
		"type": "receive"
	}
]`
	WheelSMCAddr = "0x523AC3553B4814D0a6629419ACc7adAe60aB929E"
)

func GetSignedTx() {

}

func TestEstimateGas(t *testing.T) {
	//ctx := context.Background()
	//node, err := SetupNodeClient()
	//assert.Nil(t, err)

	//gas, err := node.EstimateGas(ctx)
	//assert.Nil(t, err)
	//
	//fmt.Println("Done", gas)
}

func TestSMC_AddSpin(t *testing.T) {
	node, err := SetupNodeClient()
	assert.Nil(t, err)
	r := strings.NewReader(wheelABI)
	abiData, err := abi.JSON(r)
	assert.Nil(t, err)

	smc := NewContract(node, &abiData, common.HexToAddress(WheelSMCAddr))

	pubKey, privateKey, err := SetupTestAccount()
	assert.Nil(t, err)
	fromAddress := crypto.PubkeyToAddress(*pubKey)
	// Now we can read the nonce that we should use for the account's transaction.
	nonce, err := node.NonceAt(context.Background(), fromAddress.Hex())
	assert.Nil(t, err)
	gasLimit := uint64(3000000)
	gasPrice := big.NewInt(1)
	auth := NewKeyedTransactor(privateKey)
	auth.Nonce = nonce
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice

	spinNumber := big.NewInt(100000)

	tx, err := smc.Transact(auth, "setNumberOfSpin", common.HexToAddress("0x4f36A53DC32272b97Ae5FF511387E2741D727bdb"), spinNumber)
	assert.Nil(t, err)
	fmt.Println("TxHash", tx.Hash().String())

}

func BenchmarkWheelSpin(b *testing.B) {
	node, err := SetupNodeClient()
	assert.Nil(b, err)
	r := strings.NewReader(wheelABI)
	abiData, err := abi.JSON(r)
	assert.Nil(b, err)

	smc := NewContract(node, &abiData, common.HexToAddress(WheelSMCAddr))

	// run the Fib function b.N times
	fmt.Println("TotalRun", b.N)
	for n := 0; n < b.N; n++ {
		pubKey, privateKey, err := SetupTestAccount()
		assert.Nil(b, err)
		fromAddress := crypto.PubkeyToAddress(*pubKey)
		// Now we can read the nonce that we should use for the account's transaction.
		nonce, err := node.NonceAt(context.Background(), fromAddress.Hex())
		assert.Nil(b, err)
		gasLimit := uint64(3000000)
		gasPrice := big.NewInt(1)
		auth := NewKeyedTransactor(privateKey)
		auth.Nonce = nonce
		auth.Value = big.NewInt(0) // in wei
		auth.GasLimit = gasLimit   // in units
		auth.GasPrice = gasPrice

		_, err = smc.Transact(auth, "spin")
		assert.Nil(b, err)

	}
}

func BenchmarkFullFlow(b *testing.B) {
	fmt.Println("TotalRun", b.N)
	for i := 0; i < b.N; i++ {
		if err := FullFlow(); err != nil {
			return
		}
	}
}

func TestFullFlow(t *testing.T) {
	assert.Nil(t, FullFlow())

}

func FullFlow() error {
	ctx := context.Background()
	node, err := SetupNodeClient()
	if err != nil {
		return err
	}
	r := strings.NewReader(wheelABI)
	abiData, err := abi.JSON(r)
	if err != nil {
		return err
	}
	pubKey, privateKey, err := SetupTestAccount()
	if err != nil {
		return err
	}
	fromAddress := crypto.PubkeyToAddress(*pubKey)
	//	fmt.Println("Priv", privateKey)

	smc := NewContract(node, &abiData, common.HexToAddress(WheelSMCAddr))

	totalSpin, err := TotalSpin(ctx, node, smc, fromAddress)
	if err != nil {
		return err
	}
	//fmt.Println("TotalSpin", totalSpin)

	_, err = TotalReward(ctx, node, smc, fromAddress)
	if err != nil {
		return err
	}
	//fmt.Println("TotalReward", totalReward)

	if totalSpin > 0 {
		_, err := Spin(ctx, node, smc, fromAddress, privateKey)
		if err != nil {
			return err
		}

		_, err = TotalReward(ctx, node, smc, fromAddress)
		if err != nil {
			return err
		}
		//fmt.Println("AfterTotalReward", afterSpinTotalReward)

		//fmt.Println("Reward of spin", afterSpinTotalReward-totalReward)

	}

	return nil
}

func TotalSpin(ctx context.Context, node Node, smc *Contract, addr common.Address) (uint64, error) {
	payload, err := smc.Abi.Pack("numberOfSpin", addr)
	if err != nil {
		return 0, err
	}
	res, err := node.KardiaCall(ctx, constructCallArgs(WheelSMCAddr, payload))
	if err != nil {
		return 0, err
	}
	type reward struct {
		Reward *big.Int
	}
	var result reward

	err = smc.Abi.UnpackIntoInterface(&result, "numberOfSpin", res)
	if err != nil {
		return 0, err
	}

	return result.Reward.Uint64(), nil
}

func Spin(ctx context.Context, node Node, smc *Contract, addr common.Address, privateKey *ecdsa.PrivateKey) (string, error) {
	nonce, err := node.NonceAt(context.Background(), addr.Hex())
	if err != nil {
		return "", err
	}
	gasLimit := uint64(3000000)
	gasPrice := big.NewInt(1)
	auth := NewKeyedTransactor(privateKey)
	auth.Nonce = nonce
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice

	tx, err := smc.Transact(auth, "spin")
	if err != nil {
		return "", err
	}
	return tx.Hash().String(), nil
}

func TestSMC_WheelSpin(t *testing.T) {
	node, err := SetupNodeClient()
	assert.Nil(t, err)
	r := strings.NewReader(wheelABI)
	abiData, err := abi.JSON(r)
	assert.Nil(t, err)

	smc := NewContract(node, &abiData, common.HexToAddress(WheelSMCAddr))

	pubKey, privateKey, err := SetupTestAccount()
	assert.Nil(t, err)
	fromAddress := crypto.PubkeyToAddress(*pubKey)
	// Now we can read the nonce that we should use for the account's transaction.
	nonce, err := node.NonceAt(context.Background(), fromAddress.Hex())
	assert.Nil(t, err)
	gasLimit := uint64(3000000)
	gasPrice := big.NewInt(1)
	auth := NewKeyedTransactor(privateKey)
	auth.Nonce = nonce
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice

	tx, err := smc.Transact(auth, "spin")
	assert.Nil(t, err)
	fmt.Println("TxHash", tx.Hash().String())
}

func TotalReward(ctx context.Context, node Node, smc *Contract, fromAddr common.Address) (uint64, error) {
	payload, err := smc.Abi.Pack("reward", common.HexToAddress("0x4f36A53DC32272b97Ae5FF511387E2741D727bdb"))
	if err != nil {
		return 0, err
	}
	res, err := node.KardiaCall(ctx, constructCallArgs(WheelSMCAddr, payload))
	if err != nil {
		return 0, err
	}
	type reward struct {
		Reward *big.Int
	}
	var result reward

	if err := smc.Abi.UnpackIntoInterface(&result, "reward", res); err != nil {
		return 0, err
	}

	return result.Reward.Uint64(), nil
}

func TestSMC_WheelReward(t *testing.T) {
	ctx := context.Background()

	node, err := SetupNodeClient()
	assert.Nil(t, err)
	r := strings.NewReader(wheelABI)
	abiData, err := abi.JSON(r)
	assert.Nil(t, err)
	smc := Contract{
		Abi:             &abiData,
		ContractAddress: common.HexToAddress(WheelSMCAddr),
	}

	payload, err := smc.Abi.Pack("reward", common.HexToAddress("0x4f36A53DC32272b97Ae5FF511387E2741D727bdb"))
	assert.Nil(t, err)
	res, err := node.KardiaCall(ctx, constructCallArgs(WheelSMCAddr, payload))
	assert.Nil(t, err)
	type reward struct {
		Reward *big.Int
	}
	var result reward

	err = smc.Abi.UnpackIntoInterface(&result, "reward", res)
	assert.Nil(t, err)

	fmt.Println("Reward", result)
}
