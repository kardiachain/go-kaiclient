# KaiClient

## Initializing

-----

```go
func SetupKAIClient() (*Client, context.Context, error) {
	ctx, _ := context.WithCancel(context.Background())
	cfg := zapdriver.NewProductionConfig()
	logger, err := cfg.Build()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create logger: %v", err)
	}
	// defer logger.Sync()
	client, err := NewKaiClient("http://10.10.0.251:8551", logger)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create new KaiClient: %v", err)
	}
	return client, ctx, nil
}
```

## API

Current support API

### Blocks

------

**LatestBlockNumber**

```go
LatestBlockNumber(ctx context.Context) (uint64, error)
```

**BlockByHash**

```go
BlockByHash(ctx context.Context, hash common.Hash) (*Block, error)
```
**BlockByNumber**
```go
BlockByNumber(ctx context.Context, number uint64) (*Block, error)
```
**BlockHeaderByNumber**
```go
BlockHeaderByNumber(ctx context.Context, number uint64) (*Header, error)
```
**BlockHeaderByHash**
```go
BlockHeaderByHash(ctx context.Context, hash common.Hash) (*Header, error)
```
**GetTransaction**
```go
GetTransaction(ctx context.Context, hash common.Hash) (tx *Transaction, isPending bool, err error)
```
**GetTransactionReceipt**
```go
GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*kai.PublicReceipt, error)
```
**BalanceAt**
```go
BalanceAt(ctx context.Context, account common.Address, blockHeightOrHash interface{}) (string, error)
```
**NonceAt**
```go
NonceAt(ctx context.Context, account common.Address) (uint64, error)
```
**SendRawTransaction**
```go
SendRawTransaction(ctx context.Context, tx *Transaction) error
```
**Peers**
```go
Peers(ctx context.Context) ([]*PeerInfo, error)
```
**NodesInfo**
```go
NodesInfo(ctx context.Context) ([]*NodeInfo, error)
```
**Datadir**
```go
Datadir(ctx context.Context) (string, error)
```
**Validator**
```go
Validator(ctx context.Context, rpcURL string) (*Validator, error)
```
**Validators**
```go
Validators(ctx context.Context) ([]*Validator, error)
```
## Benchmark result