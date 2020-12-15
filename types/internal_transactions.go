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
