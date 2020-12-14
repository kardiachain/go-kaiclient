package types

type TokenHolder struct {
	TokenName       string `json:"tokenName" bson:"tokenName"`
	TokenSymbol     string `json:"tokenSymbol" bson:"tokenSymbol"`
	TokenDecimals   int64  `json:"tokenDecimals" bson:"-"`
	ContractAddress string `json:"contractAddress" bson:"contractAddress"`
	HolderAddress   string `json:"holderAddress" bson:"holderAddress"`
	Balance         string `json:"balance" bson:"balance"`

	UpdatedAt int64 `json:"updatedAt" bson:"updatedAt"`
	CreatedAt int64 `json:"createdAt" bson:"createdAt"`
}
