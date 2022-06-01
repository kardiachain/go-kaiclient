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
	"crypto/ecdsa"

	"github.com/kardiachain/go-kardia/lib/crypto"
	"github.com/kardiachain/go-kardia/rpc"
	"go.uber.org/zap"
)

var url = "https://rpc.kardiachain.io"

func setupTestNodeInterface() (Node, error) {
	lgr, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	node, err := NewNode(url, lgr)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func setupWSNodeClient() (Node, error) {
	node, err := NewNode("wss://ws-dev.kardiachain.io", zap.L())
	if err != nil {
		return nil, err
	}

	return node, nil
}

// setupTestNodeInstance return point to node struct instead interface
// for private function test
func setupTestNodeInstance() (*node, error) {
	lgr, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	node := &node{
		client: rpcClient,
		url:    url,
		lgr:    lgr,
	}
	if err := node.setupSMC(); err != nil {
		return nil, err
	}
	return node, nil
}

func setupTestAccount() (*ecdsa.PublicKey, *ecdsa.PrivateKey, error) {
	privateKey, err := crypto.HexToECDSA("63e16b5334e76d63ee94f35bd2a81c721ebbbb27e81620be6fc1c448c767eed9")
	if err != nil {
		return nil, nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, err
	}

	return publicKeyECDSA, privateKey, nil
}
