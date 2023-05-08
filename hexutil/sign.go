package hexutil

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	NullAddr = "0x0"
)

// VerifySign check if input signature is valid and return addr
func VerifySign(sign string, msg string) (bool, string) {
	loginMsgBytes := accounts.TextHash([]byte(msg))
	signBytes := hexutil.MustDecode(sign)
	signBytes[crypto.RecoveryIDOffset] -= 27

	recovered, err := crypto.SigToPub(loginMsgBytes, signBytes)
	if err != nil {
		return false, NullAddr
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)
	//sigPublicKey, err := crypto.Ecrecover(loginMsgBytes, signBytes)
	//if err != nil {
	//	return false, NullAddr
	//}
	//addr := crypto.PubkeyToAddress(*pub)
	//signatureNoRecoverID := signBytes[:len(signBytes)-1] // remove recovery id
	//verified := crypto.VerifySignature(sigPublicKey, loginMsgBytes, signatureNoRecoverID)
	return true, strings.ToLower(recoveredAddr.Hex())
}
