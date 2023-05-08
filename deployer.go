package kclient

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"

	"github.com/kardiachain/go-kaiclient/node"
)

type Deployer struct {
	Bytecode []byte
	ABI      abi.ABI

	provider *node.Eth
}

func (d *Deployer) SetProvider(p *node.Eth) error {
	d.provider = p
	return nil
}

func NewDeployer(abiString string, binString string) (*Deployer, error) {
	if len(abiString) == 0 {
		return nil, fmt.Errorf("invalid abi json string")
	}
	a, err := abi.JSON(bytes.NewReader([]byte(abiString)))
	if err != nil {
		return nil, err
	}

	deployer := &Deployer{
		Bytecode: common.FromHex(binString),
		ABI:      a,
	}
	return deployer, nil
}

func (d *Deployer) Deploy(params ...interface{}) (common.Hash, error) {
	input, err := d.ABI.Pack("", params...)
	if err != nil {
		return common.Hash{}, err
	}
	data := append(d.Bytecode, input...)
	gasLimit := uint64(3100000)
	gasPrice := big.NewInt(1000000000)
	nonce, err := d.provider.GetNonce(d.provider.Address(), nil)
	if err != nil {
		return common.Hash{}, err
	}

	txHash, err := d.provider.SendDeployTransaction(
		big.NewInt(0),
		nonce,
		gasLimit,
		gasPrice,
		data,
	)
	fmt.Println("TxHash", txHash.String())
	return txHash, nil
}
