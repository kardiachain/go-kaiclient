package node

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/kardiachain/go-kaiclient/types"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/crypto"
	"github.com/kardiachain/go-kardia/lib/rlp"
	eTypes "github.com/kardiachain/go-kardia/types"
)

// SignTx signs the transaction using the given signer and private key.
func SignTx(tx *eTypes.Transaction, s eTypes.Signer, prv *ecdsa.PrivateKey) (*eTypes.Transaction, error) {
	h := s.Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(s, sig)
}

func (e *Eth) SendRawTransaction(
	to common.Address,
	amount *big.Int,
	nonce uint64,
	gasLimit uint64,
	gasPrice *big.Int,
	data []byte,
) (common.Hash, error) {
	var hash common.Hash

	tx := eTypes.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)

	signedTx, err := eTypes.SignTx(eTypes.HomesteadSigner{}, tx, e.privateKey)
	if err != nil {
		return hash, err
	}
	serializedTx, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return hash, err
	}

	err = e.c.Call("eth_sendRawTransaction", &hash, fmt.Sprintf("0x%x", serializedTx))
	return hash, err

}

func (e *Eth) SyncSendRawTransaction(
	to common.Address,
	amount *big.Int,
	nonce uint64,
	gasLimit uint64,
	gasPrice *big.Int,
	data []byte,
) (*types.Receipt, error) {

	tx := eTypes.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)

	signedTx, err := eTypes.SignTx(eTypes.HomesteadSigner{}, tx, e.privateKey)
	if err != nil {
		return nil, err
	}
	serializedTx, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return nil, err
	}
	var hash common.Hash
	err = e.c.Call("eth_sendRawTransaction", &hash, fmt.Sprintf("0x%x", serializedTx))
	if err != nil {
		return nil, err
	}

	// fmt.Printf("hash %v\n", hash)

	type ReceiptCh struct {
		ret *types.Receipt
		err error
	}

	var timeoutFlag int32
	ch := make(chan *ReceiptCh, 1)

	go func() {
		for {
			receipt, _ := e.GetTransactionReceipt(hash)
			// if err != nil && err.Error() != "not found" {
			// 	ch <- &ReceiptCh{
			// 		err: err,
			// 	}
			// 	break
			// }
			if receipt != nil {
				ch <- &ReceiptCh{
					ret: receipt,
					err: nil,
				}
				break
			}
			if atomic.LoadInt32(&timeoutFlag) == 1 {
				break
			}
		}
		// fmt.Println("send tx done")
	}()

	select {
	case result := <-ch:
		return result.ret, result.err
	case <-time.After(time.Duration(e.txPollTimeout) * time.Second):
		atomic.StoreInt32(&timeoutFlag, 1)
		return nil, fmt.Errorf("transaction was not mined within %v seconds, "+
			"please make sure your transaction was properly sent. Be aware that it might still be mined", e.txPollTimeout)
	}
}

func (e *Eth) SendDeployTransaction(
	amount *big.Int,
	nonce uint64,
	gasLimit uint64,
	gasPrice *big.Int,
	data []byte,
) (common.Hash, error) {
	var hash common.Hash

	tx := eTypes.NewContractCreation(nonce, amount, gasLimit, gasPrice, data)

	signedTx, err := eTypes.SignTx(eTypes.HomesteadSigner{}, tx, e.privateKey)
	if err != nil {
		return hash, err
	}
	serializedTx, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return hash, err
	}

	err = e.c.Call("eth_sendRawTransaction", &hash, fmt.Sprintf("0x%x", serializedTx))
	return hash, err

}
