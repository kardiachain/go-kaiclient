/*
 *  Copyright 2018 KardiaChain
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

func (p *Provider) GetInsertBlockTime() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.insertBlockTime.String()
}

func (p *Provider) GetRawInsertBlockTime() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.insertBlockTime.Raw()
}

func (p *Provider) GetScrapingTime() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.scrapingTime.String()
}

func (p *Provider) GetRawScrapingTime() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.scrapingTime.Raw()
}

func (p *Provider) GetInsertTxsTime() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.insertTxsTime.String()
}

func (p *Provider) GetRawInsertTxsTime() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.insertTxsTime.Raw()
}

func (p *Provider) GetInsertActiveAddressTime() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.insertActiveAddressTime.String()
}

func (p *Provider) GetRawInsertActiveAddressTime() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.insertActiveAddressTime.Raw()
}

func (p *Provider) GetUpsertBlockTime() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.upsertBlockTime.String()
}

func (p *Provider) GetRawUpsertBlockTime() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.upsertBlockTime.Raw()
}

func (p *Provider) GetLatestBLock() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.latestBlock
}

func (p *Provider) GetTodoLength() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.todoLength
}

func (p *Provider) GetReorgedBlocks() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.reorgedBlocks
}

func (p *Provider) GetInvalidBlocks() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.invalidBlocks
}
