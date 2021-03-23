// Package kardia
package kardia

import (
	"context"
	"math/big"

	"github.com/kardiachain/go-kardia/lib/common"
)

type token struct {
	node Node
	c    *Contract
}

func (t *token) IsKRC20(ctx context.Context) bool {
	if t.c.ContractAddress.Equal(common.Address{}) {
		return false
	}
	if t, err := t.KRC20Info(ctx, t.c); err != nil || t == nil {
		return false
	}

	return true
}

func NewToken(node Node, c *Contract) Token {
	return &token{
		node: node,
		c:    c,
	}
}

type Token interface {
	IsKRC20(ctx context.Context) bool
	KRC20Info(ctx context.Context, c *Contract) (*KRC20, error)
	HolderBalance(ctx context.Context, c *Contract, holderAddress common.Address) (*big.Int, error)
}

func (t *token) KRC20Info(ctx context.Context, c *Contract) (*KRC20, error) {
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

func (t *token) HolderBalance(ctx context.Context, c *Contract, holderAddress common.Address) (*big.Int, error) {
	return nil, nil
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