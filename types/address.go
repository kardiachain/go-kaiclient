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

package types

type Address struct {
	Address string `json:"address" bson:"address"`
	Balance string `json:"balance" bson:"balance"`
	Name    string `json:"name" bson:"name"` // alias of an address

	// Token
	TokenName   string `json:"tokenName" bson:"tokenName"`
	TokenSymbol string `json:"tokenSymbol" bson:"tokenSymbol"`
	Decimals    int64  `json:"decimals" bson:"decimals"`
	TotalSupply string `json:"totalSupply" bson:"totalSupply"`

	// SMC
	IsContract   bool     `json:"isContract" bson:"isContract"`
	ErcTypes     []string `json:"ercTypes" bson:"ercTypes"`
	Interfaces   []string `json:"interfaces" bson:"interfaces"`
	OwnerAddress string   `json:"ownerAddress" bson:"ownerAddress"`

	// Stats
	TxCount         int `json:"txCount" bson:"txCount"`
	HolderCount     int `json:"holderCount" bson:"holderCount"`
	InternalTxCount int `json:"internalTxCount" bson:"internalTxCount"`
	TokenTxCount    int `json:"tokenTxCount" bson:"tokenTxCount"`

	UpdatedAt int64 `json:"updatedAt" bson:"updatedAt"`
}
