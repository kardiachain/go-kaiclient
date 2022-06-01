// Package kardia
package kardia

import (
	"go.uber.org/zap"
	"strings"

	"github.com/kardiachain/go-kaiclient/kardia/smc"
	"github.com/kardiachain/go-kardia/lib/abi"
)

func KRC721ABI() (*abi.ABI, error) {
	lgr, _ := zap.NewDevelopment()
	lgr.Info("KRC721 ABI", zap.String("ABI", smc.KRC721ABI))

	r := strings.NewReader(smc.KRC721ABI)
	abiData, err := abi.JSON(r)
	if err != nil {
		return nil, err
	}
	return &abiData, nil
}

func KRC20ABI() (*abi.ABI, error) {
	r := strings.NewReader(smc.KRC20ABI)
	abiData, err := abi.JSON(r)
	if err != nil {
		return nil, err
	}
	return &abiData, nil
}
