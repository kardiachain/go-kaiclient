# KaiClient


## Initializing

-----

```go
func SetupKardiaClient() (Node, error){
	url := "https://dev-1.kardiachain.io"
	lgr, err := zap.NewProduction()
	if err != nil {
        return nil, err
    }
    node, err := NewNode(url, lgr)
    if err != nil {
    	return nil, err
    }
    return node, nil
}

```

## API

### Info

------

```go
type IInfo interface {
    Url() string
    IsAlive() bool
    NodeInfo(ctx context.Context) (*NodeInfo, error)
    GetCirculatingSupply(ctx context.Context) (*big.Int, error)
    KardiaCall(ctx context.Context, args SMCCallArgs) ([]byte, error)
}
```

### Blocks

------

```go
type IBlock interface {
	LatestBlockNumber(ctx context.Context) (uint64, error)
	BlockByHash(ctx context.Context, hash string) (*Block, error)
	BlockByHeight(ctx context.Context, height uint64) (*Block, error)
	BlockHeaderByHash(ctx context.Context, hash string) (*Header, error)
	BlockHeaderByNumber(ctx context.Context, number uint64) (*Header, error)
}
```

### Addresses

------

```go
type IAddress interface {
    Balance(ctx context.Context, addressHash string) (string, error)
    StorageAt(ctx context.Context, addressHash string, key string) ([]byte, error)
    Code(ctx context.Context, addressHash string) (common.Bytes, error)
    NonceAt(ctx context.Context, addressHash string) (uint64, error)
}
```

### Transactions

------

```go
type ITx interface {
    GetTransaction(ctx context.Context, hash string) (*Transaction, error)
    GetTransactionReceipt(ctx context.Context, txHash string) (*Receipt, error)
    SendTransaction(ctx context.Context, tx *types.Transaction) error
    SendRawTransaction(ctx context.Context, tx *types.Transaction) error
}
```

## Examples

_Note:_ Examples can be found at *_test.go

---

### Interact with SMC

Function definition

```shell
{
    "constant": true,
    "inputs": [
      {
        "internalType": "address",
        "name": "_delAddr",
        "type": "address"
      }
    ],
    "name": "getDelegationRewards",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
```

#### Create node instance

```go
n, err := NewNode("https://dev-1.kardiachain.io", zap.L())
if err != nil {
	return 
}

```

#### Create contract instance

```go
// validatorSMC ABI can be found at smc package
validatorSmcAbi, err := abi.JSON(strings.NewReader(smc.ValidatorABI))
if err != nil {
    return err
}
validatorUtil := &Contract{
    Abi: &validatorSmcAbi,
}
n.validatorSMC = validatorUtil
```

#### Build payload

```go
payload, err := n.validatorSMC.Abi.Pack("getDelegationRewards", common.HexToAddress(delegatorAddress))
if err != nil {
    n.lgr.Error("Error packing delegation rewards payload: ", zap.Error(err))
    return nil, err
}
```

#### Send payload

```go
res, err := n.KardiaCall(ctx, ConstructCallArgs(validatorSMCAddr, payload))
if err != nil {
    n.lgr.Error("GetDelegationRewards KardiaCall error: ", zap.Error(err))
    return nil, err
}
```

#### Get response

```go
var result struct {
    Rewards *big.Int
}
// unpack result
err = n.validatorSMC.Abi.UnpackIntoInterface(&result, "getDelegationRewards", res)
if err != nil {
    n.lgr.Error("Error unpacking delegation rewards: ", zap.Error(err))
    return nil, err
}
```

---

### Subscribe event

```go
func TestSubscription_NewBlockHead(t *testing.T) {
	lgr, err := zap.NewProduction()
	assert.Nil(t, err)
	url := "wss://ws-dev.kardiachain.io/ws"

	node, err := NewNode(url, lgr)
	assert.Nil(t, err)

	headersCh := make(chan *types.Header)
	sub, err := node.SubscribeNewHead(context.Background(), headersCh)
	assert.Nil(t, err, "cannot subscribe")

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headersCh:
			fmt.Println(header.Hash().Hex())
		}
	}
}

```