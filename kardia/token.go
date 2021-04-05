// Package kardia
package kardia

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"

	"github.com/kardiachain/go-kaiclient/kardia/smc"
)

type token struct {
	node Node
	c    *Contract
}

func (t *token) isKRC20() bool {
	if t.c.ContractAddress.Equal(common.Address{}) {
		return false
	}
	if t, err := t.KRC20Info(context.Background()); err != nil || t == nil {
		return false
	}

	return true
}

func NewKRC20(node Node, address string, owner string) (Token, error) {
	r := strings.NewReader(smc.KRC20ABI)
	abiData, err := abi.JSON(r)
	if err != nil {
		return nil, err
	}
	c := &Contract{
		Abi:             &abiData,
		Bytecode:        smc.KRC20Bytecode,
		ContractAddress: common.HexToAddress(address),
		OwnerAddress:    common.HexToAddress(owner),
	}

	t := &token{
		node: node,
		c:    c,
	}
	if !t.isKRC20() {
		return nil, fmt.Errorf("not valid KRC20")
	}

	return &token{
		node: node,
		c:    c,
	}, nil
}

type Token interface {
	KRC20Info(ctx context.Context) (*KRC20, error)
	HolderBalance(ctx context.Context, holderAddress common.Address) (*big.Int, error)
}

func (t *token) KRC20Info(ctx context.Context) (*KRC20, error) {
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

func (t *token) HolderBalance(ctx context.Context, holderAddress common.Address) (*big.Int, error) {
	payload, err := t.c.Abi.Pack("balanceOf", holderAddress)
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
