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

import (
	"time"
)

const (
	defaultLimit = 50
	MaximumLimit = 100
)

type Pagination struct {
	Skip  int
	Limit int
}

func (f *Pagination) Sanitize() {
	if f.Skip < 0 {
		f.Skip = 0
	}
	if f.Limit <= 0 {
		f.Limit = defaultLimit
	} else if f.Limit > MaximumLimit {
		f.Limit = MaximumLimit
	}
}

type SortFilter struct {
	SortBy string
	Asc    bool
}

type TimeFilter struct {
	FromTime time.Time
	ToTime   time.Time
}

func (f *TimeFilter) Sanitize() {
	if f.FromTime.IsZero() {
		f.FromTime = time.Unix(0, 0)
	}
	if f.ToTime.IsZero() {
		f.ToTime = time.Now()
	}
}

type ContractsFilter struct {
	Pagination
	SortFilter
	ContractName string
	TokenName    string
	TokenSymbol  string
	ErcType      string
}

type InternalTxsFilter struct {
	Pagination
	TokenTransactions bool
	InternalAddress   string
}

type TxsFilter struct {
	Pagination
	TimeFilter
}

type BlocksFilter struct {
	Pagination
	TimeFilter
}

func (f *TxsFilter) Sanitize() {
	f.Pagination.Sanitize()
	f.TimeFilter.Sanitize()
}
