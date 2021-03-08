// Package kardia
package kardia

import (
	"crypto/ecdsa"

	"github.com/kardiachain/go-kardia/lib/crypto"
	"go.uber.org/zap"
)

var url = "https://rpc.kardiachain.io"

func SetupNodeClient() (Node, error) {
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

func SetupWSNodeClient() (Node, error) {
	node, err := NewNode("ws://10.10.0.251:8550/ws", zap.L())
	if err != nil {
		return nil, err
	}

	return node, nil
}

func SetupDevNodeClient() (Node, error) {
	node, err := NewNode("https://dev-1.kardiachain.io", zap.L())
	if err != nil {
		return nil, err
	}

	return node, nil
}

func SetupTestAccount() (*ecdsa.PublicKey, *ecdsa.PrivateKey, error) {
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
