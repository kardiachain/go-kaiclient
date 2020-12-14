package types

// TokenTransfer represents a Transfer event emitted from an ERC20 or ERC721.
type TokenTransfer struct {
	TransactionHash string `json:"transactionHash" bson:"transactionHash"`
	Contract        string `json:"contractAddress" bson:"contractAddress"`

	From        string `json:"fromAddress" bson:"fromAddress"`
	To          string `json:"toAddress" bson:"toAddress"`
	Value       string `json:"value" bson:"value"`
	BlockNumber uint64 `json:"blockNumber" bson:"blockNumber"`

	UpdatedAt int64 `json:"updatedAt" bson:"updatedAt"`
}
