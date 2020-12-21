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
package metrics

import (
	"sync"
	"time"
)

type Provider struct {
	mu sync.Mutex

	insertBlockTime         AverageDuration
	scrapingTime            AverageDuration
	insertTxsTime           AverageDuration
	insertActiveAddressTime AverageDuration
	upsertBlockTime         AverageDuration

	latestBlock   int64
	todoLength    int64
	reorgedBlocks int64
	invalidBlocks int64
}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Reset() {
	p.insertBlockTime.Reset()
	p.scrapingTime.Reset()
	p.insertTxsTime.Reset()
	p.insertActiveAddressTime.Reset()

	p.latestBlock = 0
	p.todoLength = 0
	p.reorgedBlocks = 0
	p.invalidBlocks = 0
}

func (p *Provider) RecordInsertBlockTime(duration time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.insertBlockTime.Add(duration)
}

func (p *Provider) RecordScrapingTime(duration time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.scrapingTime.Add(duration)
}

func (p *Provider) RecordInsertTxsTime(duration time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.insertTxsTime.Add(duration)
}

func (p *Provider) RecordInsertActiveAddressTime(duration time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.insertActiveAddressTime.Add(duration)
}

func (p *Provider) RecordUpsertBlockTime(duration time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.upsertBlockTime.Add(duration)
}

func (p *Provider) RecordLatestBlock(block int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.latestBlock = block
}

func (p *Provider) RecordTodoLength(len int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.todoLength = len
}

func (p *Provider) RecordReorgedBlock() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.reorgedBlocks++
}

func (p *Provider) RecordInvalidBlock() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.invalidBlocks++
}
