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
	"math/big"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
)

const (
	TokenTypeUnknown = 0
	TokenTypeKRC20   = 1
	TokenTypeKRC721  = 2
)

type Token interface {
	TokenType() int
	KRC20Info(ctx context.Context) (*KRC20, error)
	KRC721Info(ctx context.Context) (*KRC721, error)
	HolderBalance(ctx context.Context, holderAddress string) (*big.Int, error)
	TotalSupply(ctx context.Context) (*big.Int, error)
	ABI() *abi.ABI
}

type token struct {
	node    Node
	c       *Contract
	krcType int
}

func NewToken(node Node, address string) (Token, error) {
	c := &Contract{
		ContractAddress: common.HexToAddress(address),
	}

	t := &token{
		node: node,
		c:    c,
	}

	//krcType := t.getKRCType(ctx)
	//if krcType == TokenTypeUnknown {
	//	return nil, errors.New("token type not support")
	//}
	//t.krcType = krcType
	return t, nil
}

func (t *token) TokenType() int {
	return t.krcType
}

func (t *token) getKRCType(ctx context.Context) int {
	if krc20Info, err := t.KRC20Info(ctx); err == nil && krc20Info != nil {
		return TokenTypeKRC20
	}
	if krc721Info, err := t.KRC721Info(ctx); err == nil && krc721Info != nil {
		return TokenTypeKRC721
	}

	return TokenTypeUnknown
}

func (t *token) KRC721Info(ctx context.Context) (*KRC721, error) {
	krc721ABI, err := KRC721ABI()
	if err != nil {
		return nil, err
	}
	t.c.Abi = krc721ABI
	name, err := t.getName(ctx)
	if err != nil {
		return nil, err
	}

	symbol, err := t.getSymbol(ctx)
	if err != nil {
		return nil, err
	}

	totalSupply, err := t.getTotalSupply(ctx)
	if err != nil {
		return nil, err
	}
	krc721 := &KRC721{
		Address:     t.c.ContractAddress,
		Name:        name,
		Symbol:      symbol,
		TotalSupply: totalSupply,
	}
	return krc721, nil
}

func (t *token) KRC20Info(ctx context.Context) (*KRC20, error) {
	krc20ABI, err := KRC20ABI()
	if err != nil {
		return nil, err
	}
	t.c.Abi = krc20ABI
	name, err := t.getName(ctx)
	if err != nil {
		return nil, err
	}

	symbol, err := t.getSymbol(ctx)
	if err != nil {
		return nil, err
	}

	decimals, err := t.getDecimals(ctx)
	if err != nil {
		return nil, err
	}

	totalSupply, err := t.getTotalSupply(ctx)
	if err != nil {
		return nil, err
	}

	krc20 := &KRC20{
		Address:     t.c.ContractAddress,
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
	}
	return krc20, nil
}

func (t *token) ABI() *abi.ABI {
	return t.c.ABI()
}

func (t *token) HolderBalance(ctx context.Context, holderAddress string) (*big.Int, error) {
	address := common.HexToAddress(holderAddress)
	payload, err := t.c.Abi.Pack("balanceOf", address)
	if err != nil {
		return nil, err
	}

	res, err := t.node.KardiaCall(ctx, ConstructCallArgs(t.c.ContractAddress.Hex(), payload))
	if err != nil {
		return nil, err
	}

	var balance *big.Int
	// unpack result
	err = t.c.Abi.UnpackIntoInterface(&balance, "balanceOf", res)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (t *token) TotalSupply(ctx context.Context) (*big.Int, error) {
	return t.getTotalSupply(ctx)
}

func (t *token) getName(ctx context.Context) (string, error) {
	payload, err := t.c.Abi.Pack("name")
	if err != nil {
		return "", err
	}
	res, err := t.node.KardiaCall(ctx, ConstructCallArgs(t.c.ContractAddress.String(), payload))
	if err != nil {
		return "", err
	}

	var name string

	if err := t.c.Abi.UnpackIntoInterface(&name, "name", res); err != nil {
		return "", err
	}

	return name, nil

}

func (t *token) getSymbol(ctx context.Context) (string, error) {
	payload, err := t.c.Abi.Pack("symbol")
	if err != nil {
		return "", err
	}
	res, err := t.node.KardiaCall(ctx, ConstructCallArgs(t.c.ContractAddress.String(), payload))
	if err != nil {
		return "", err
	}
	var symbol string

	if err := t.c.Abi.UnpackIntoInterface(&symbol, "symbol", res); err != nil {
		return "", err
	}

	return symbol, nil
}

func (t *token) getDecimals(ctx context.Context) (uint8, error) {
	payload, err := t.c.Abi.Pack("decimals")
	if err != nil {
		return 0, err
	}
	res, err := t.node.KardiaCall(ctx, ConstructCallArgs(t.c.ContractAddress.String(), payload))
	if err != nil {
		return 0, err
	}
	var decimals uint8

	if err := t.c.Abi.UnpackIntoInterface(&decimals, "decimals", res); err != nil {
		return 0, err
	}

	return decimals, nil
}

func (t *token) getTotalSupply(ctx context.Context) (*big.Int, error) {
	payload, err := t.c.Abi.Pack("totalSupply")
	if err != nil {
		return nil, err
	}
	res, err := t.node.KardiaCall(ctx, ConstructCallArgs(t.c.ContractAddress.String(), payload))
	if err != nil {
		return nil, err
	}
	var totalSupply *big.Int

	if err := t.c.Abi.UnpackIntoInterface(&totalSupply, "totalSupply", res); err != nil {
		return nil, err
	}

	return totalSupply, nil
}

func (t *token) getOwnerBalance(ctx context.Context) (*big.Int, error) {
	payload, err := t.c.Abi.Pack("balanceOf", t.c.OwnerAddress)
	if err != nil {
		return nil, err
	}
	res, err := t.node.KardiaCall(ctx, ConstructCallArgs(t.c.ContractAddress.String(), payload))
	if err != nil {
		return nil, err
	}

	var balance *big.Int

	if err := t.c.Abi.UnpackIntoInterface(&balance, "balanceOf", res); err != nil {
		return nil, err
	}

	return balance, nil
}
