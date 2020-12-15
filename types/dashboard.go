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
// Package types
package types

type TokenInfo struct {
	Name              string  `json:"name"`
	Symbol            string  `json:"symbol"`
	Decimal           int64   `json:"decimal"`
	TotalSupply       int64   `json:"total_supply"`
	CirculatingSupply int64   `json:"circulating_supply"`
	Price             float64 `json:"price"`
	Volume24h         float64 `json:"volume_24h"`
	Change1h          float64 `json:"change_1h"`
	Change24h         float64 `json:"change_24h"`
	Change7d          float64 `json:"change_7d"`
	MarketCap         float64 `json:"market_cap"`
}
