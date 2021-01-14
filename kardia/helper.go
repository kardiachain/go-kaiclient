// Package kardia
package kardia

import (
	"context"
	"math/big"

	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/rpc"
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
	err = client.CallContext(ctx, &res, "kai_kardiaCall", constructCallArgs(stakingSMC.ContractAddress.Hex(), payload), "latest")
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

func constructCallArgs(address string, payload []byte) SMCCallArgs {
	return SMCCallArgs{
		From:     address,
		To:       &address,
		Gas:      100000000,
		GasPrice: big.NewInt(0),
		Value:    big.NewInt(0),
		Data:     common.Bytes(payload).String(),
	}
}
