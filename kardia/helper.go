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
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/abi/bind"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/crypto"
	"github.com/kardiachain/go-kardia/rpc"
	"github.com/kardiachain/go-kardia/types"
	"github.com/shopspring/decimal"
)

func getParamsSMCAddress(stakingSMC *Contract, client *rpc.Client) (common.Address, error) {
	payload, err := stakingSMC.Abi.Pack("params")
	if err != nil {
		return common.Address{}, err
	}

	var (
		res common.Bytes
		ctx = context.Background()
	)
	err = client.CallContext(ctx, &res, "kai_kardiaCall", ConstructCallArgs(stakingSMC.ContractAddress.Hex(), payload), "latest")
	if err != nil {
		return common.Address{}, err
	}

	var result struct {
		ParamsSmcAddr common.Address
	}
	err = stakingSMC.Abi.UnpackIntoInterface(&result, "params", res)
	if err != nil {
		return common.Address{}, err
	}

	return result.ParamsSmcAddr, nil
}

func ConstructCallArgs(address string, payload []byte) SMCCallArgs {
	return SMCCallArgs{
		From:     address,
		To:       &address,
		Gas:      100000000,
		GasPrice: big.NewInt(0),
		Value:    big.NewInt(0),
		Data:     common.Bytes(payload).String(),
	}
}

// NewKeyedTransactor is a utility method to easily create a transaction signer
// from a single private key.
func NewKeyedTransactor(key *ecdsa.PrivateKey) *bind.TransactOpts {
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	return &bind.TransactOpts{
		From: keyAddr,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}
			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
}

func DecodeWithABI(input string, a *abi.ABI) (*FunctionCall, error) {
	// return nil if input data is too short
	if len(input) <= 2 {
		return nil, nil
	}
	data, err := hex.DecodeString(strings.TrimLeft(input, "0x"))
	if err != nil {
		return nil, err
	}
	sig := data[0:4] // get the function signature (first 4 bytes of input data)
	method, err := a.MethodById(sig)
	if err != nil {
		return nil, err
	}

	fmt.Println("Method", method)
	// exclude the function signature, only decode and unpack the arguments
	var body []byte
	if len(data) <= 4 {
		body = []byte{}
	} else {
		body = data[4:]
	}
	args, err := getInputArguments(a, method.Name, body)
	if err != nil {
		return nil, err
	}
	arguments := make(map[string]interface{})
	err = args.UnpackIntoMap(arguments, body)
	if err != nil {
		return nil, err
	}
	// convert address, bytes and string arguments into their hex representations
	for i, arg := range arguments {
		arguments[i] = parseBytesArrayIntoString(arg)
	}
	return &FunctionCall{
		Function:   method.String(),
		MethodID:   "0x" + hex.EncodeToString(sig),
		MethodName: method.Name,
		Arguments:  arguments,
	}, nil
}

func getInputArguments(a *abi.ABI, name string, data []byte) (abi.Arguments, error) {
	var args abi.Arguments
	if method, ok := a.Methods[name]; ok {
		if len(data)%32 != 0 {
			return nil, fmt.Errorf("abi: improperly formatted output: %s - Bytes: [%+v]", string(data), data)
		}
		args = method.Inputs
	}
	if args == nil {
		return nil, ErrMethodNotFound
	}
	return args, nil
}

func CheckAddress(address string) string {
	return common.HexToAddress(address).String()
}

func GenerateWallet() (common.Address, ecdsa.PrivateKey, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return common.Address{}, ecdsa.PrivateKey{}, err
	}
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, ecdsa.PrivateKey{}, fmt.Errorf("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return address, *privKey, nil
}

func FloatToBigInt(amount float64, decimals int32) *big.Int {
	return decimal.NewFromFloat(amount).Mul(decimal.New(1, decimals)).BigInt()
}
