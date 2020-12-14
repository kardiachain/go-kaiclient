# KaiClient
## Initializing
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
## Endpoints
**LatestBlockNumber**
```go
LatestBlockNumber(ctx context.Context) (uint64, error)
```
**BlockByHash**
```go
BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
```
**BlockByNumber**
```go
BlockByNumber(ctx context.Context, number uint64) (*types.Block, error)
```
**BlockHeaderByNumber**
```go
BlockHeaderByNumber(ctx context.Context, number uint64) (*types.Header, error)
```
**BlockHeaderByHash**
```go
BlockHeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
```
**GetTransaction**
```go
GetTransaction(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
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
SendRawTransaction(ctx context.Context, tx *types.Transaction) error
```
**Peers**
```go
Peers(ctx context.Context) ([]*types.PeerInfo, error)
```
**NodesInfo**
```go
NodesInfo(ctx context.Context) ([]*types.NodeInfo, error)
```
**Datadir**
```go
Datadir(ctx context.Context) (string, error)
```
**Validator**
```go
Validator(ctx context.Context, rpcURL string) (*types.Validator, error)
```
**Validators**
```go
Validators(ctx context.Context) ([]*types.Validator, error)
```
## Benchmark result