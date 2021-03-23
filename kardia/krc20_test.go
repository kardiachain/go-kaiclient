// Package kardia
package kardia

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/stretchr/testify/assert"

	"github.com/kardiachain/go-kaiclient/kardia/smc"
)

func TestValidator_KRC20(t *testing.T) {
	node, err := SetupNodeClient()
	assert.Nil(t, err)
	r := strings.NewReader(smc.KRC20ABI)
	abiData, err := abi.JSON(r)
	assert.Nil(t, err)
	c := &Contract{
		Abi:             &abiData,
		Bytecode:        smc.KRC20Bytecode,
		ContractAddress: common.HexToAddress("0x5a41FCdc110F7C39C02ea7AFc033Afa90362BC64"),
		OwnerAddress:    common.HexToAddress("0x4f36A53DC32272b97Ae5FF511387E2741D727bdb"),
	}
	validator := &ContractValidator{
		node: node,
		c:    c,
	}
	ctx := context.Background()
	name, err := validator.getName(ctx)
	assert.Nil(t, err)
	fmt.Println("Name", name)
	symbol, err := validator.getSymbol(ctx)
	assert.Nil(t, err)
	fmt.Println("Symbol", symbol)

	decimal, err := validator.getDecimals(ctx)
	assert.Nil(t, err)
	fmt.Println("Decimals", decimal)
	supply, err := validator.getTotalSupply(ctx)
	assert.Nil(t, err)
	fmt.Println("TotalSupply", supply)

	balance, err := validator.getOwnerBalance(ctx)
	assert.Nil(t, err)
	fmt.Println("Balance", balance.String())
	//name, err := validator.getName(ctx)
	//assert.Nil(t, err)
	//name, err := validator.getName(ctx)
	//assert.Nil(t, err)

}
