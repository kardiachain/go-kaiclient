// Package kardia
package kardia

import (
	"context"
	"math/big"
)

// KRC20 contract basic info
type KRC20 struct {
	node Node
	c    *Contract
}

func (k *KRC20) getName(ctx context.Context) (string, error) {
	payload, err := k.c.Abi.Pack("name")
	if err != nil {
		return "", err
	}
	res, err := k.node.KardiaCall(ctx, ConstructCallArgs(k.c.ContractAddress.String(), payload))
	if err != nil {
		return "", err
	}

	var name string

	if err := k.c.Abi.UnpackIntoInterface(&name, "name", res); err != nil {
		return "", err
	}

	return name, nil

}

func (k *KRC20) getSymbol(ctx context.Context) (string, error) {
	payload, err := k.c.Abi.Pack("symbol")
	if err != nil {
		return "", err
	}
	res, err := k.node.KardiaCall(ctx, ConstructCallArgs(k.c.ContractAddress.String(), payload))
	if err != nil {
		return "", err
	}
	var symbol string

	if err := k.c.Abi.UnpackIntoInterface(&symbol, "symbol", res); err != nil {
		return "", err
	}

	return symbol, nil
}

func (k *KRC20) getDecimals(ctx context.Context) (uint8, error) {
	payload, err := k.c.Abi.Pack("decimals")
	if err != nil {
		return 0, err
	}
	res, err := k.node.KardiaCall(ctx, ConstructCallArgs(k.c.ContractAddress.String(), payload))
	if err != nil {
		return 0, err
	}
	var decimals uint8

	if err := k.c.Abi.UnpackIntoInterface(&decimals, "decimals", res); err != nil {
		return 0, err
	}

	return decimals, nil
}

func (k *KRC20) getTotalSupply(ctx context.Context) (*big.Int, error) {
	payload, err := k.c.Abi.Pack("totalSupply")
	if err != nil {
		return nil, err
	}
	res, err := k.node.KardiaCall(ctx, ConstructCallArgs(k.c.ContractAddress.String(), payload))
	if err != nil {
		return nil, err
	}
	var totalSupply *big.Int

	if err := k.c.Abi.UnpackIntoInterface(&totalSupply, "totalSupply", res); err != nil {
		return nil, err
	}

	return totalSupply, nil
}

func (k *KRC20) getOwnerBalance(ctx context.Context) (*big.Int, error) {
	payload, err := k.c.Abi.Pack("balanceOf", k.c.OwnerAddress)
	if err != nil {
		return nil, err
	}
	res, err := k.node.KardiaCall(ctx, ConstructCallArgs(k.c.ContractAddress.String(), payload))
	if err != nil {
		return nil, err
	}

	var balance *big.Int

	if err := k.c.Abi.UnpackIntoInterface(&balance, "balanceOf", res); err != nil {
		return nil, err
	}

	return balance, nil
}
