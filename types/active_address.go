package types

type ActiveAddress struct {
	Address    string `json:"address" bson:"address"`
	Balance    string `json:"balance" bson:"balance"`
	IsContract bool   `json:"isContract" bson:"isContract"`
}
