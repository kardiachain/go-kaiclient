// Package kardia
package kardia

import (
	"context"
	"fmt"

	"github.com/kardiachain/go-kardia/lib/common"
)

type IToken interface {
	TokenInfo(ctx context.Context, c *Contract) (*Token, error)
}

func (n *node) TokenInfo(ctx context.Context, c *Contract) (*Token, error) {
	k := KRC20{
		node: n,
		c:    c,
	}
	name, err := k.getName(ctx)
	if err != nil {
		return nil, err
	}

	symbol, err := k.getSymbol(ctx)
	if err != nil {
		return nil, err
	}

	decimals, err := k.getDecimals(ctx)
	if err != nil {
		return nil, err
	}

	totalSupply, err := k.getTotalSupply(ctx)
	if err != nil {
		return nil, err
	}

	token := &Token{
		Address:     k.c.ContractAddress,
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
	}
	return token, nil
}

type TokenValidator interface {
	IsKRC20(ctx context.Context, c *Contract) bool
	IsKRC721(ctx context.Context, c *Contract) bool
}

type tokenValidator struct {
	node  Node
	krc20 KRC20
}

func NewTokenValidator(node Node) *tokenValidator {
	return &tokenValidator{
		node: node,
	}
}

//IsKRC20 check if contract implement KRC20 interface
//Current version is really simple
func (v *tokenValidator) IsKRC20(ctx context.Context, c *Contract) (bool, error) {
	if c.ContractAddress.Equal(common.Address{}) {
		return false, fmt.Errorf("contract address must not empty")
	}

	if c.OwnerAddress.Equal(common.Address{}) {
		return false, fmt.Errorf("owner address must not empty")
	}

	if c.Abi == nil {
		return false, fmt.Errorf("abi must not nil")
	}

	// Create KRC20 instance and call
	krc20 := KRC20{
		node: v.node,
		c:    c,
	}
	if _, err := krc20.getName(ctx); err != nil {
		return false, nil
	}

	if _, err := krc20.getSymbol(ctx); err != nil {
		return false, nil
	}

	if _, err := krc20.getDecimals(ctx); err != nil {
		return false, nil
	}

	if _, err := krc20.getTotalSupply(ctx); err != nil {
		return false, nil
	}

	if _, err := krc20.getOwnerBalance(ctx); err != nil {
		return false, nil
	}
	return true, nil
}

func (v *tokenValidator) IsKRC721() (bool, error) {
	return false, nil
}
