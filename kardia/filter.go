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
// Package kardia
package kardia

import (
	"encoding/hex"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
)

type Filter interface {
	Events(log Log) ([]*Event, error)
}

type filter struct {
	subscribedEvents map[string]bool
	abi              *abi.ABI
}

func NewFilter(events []string, abi *abi.ABI) (Filter, error) {
	filter := &filter{
		abi: abi,
	}
	filter.subscribedEvents = make(map[string]bool)
	for _, ev := range events {
		filter.subscribedEvents[ev] = true
	}
	return filter, nil
}

func (f *filter) Events(log Log) ([]*Event, error) {
	var events []*Event
	for _, t := range log.Topics {
		nEvent, err := f.abi.EventByID(common.HexToHash(t))
		if err != nil {
			continue
		}
		sub, found := f.subscribedEvents[nEvent.RawName]
		if !found || !sub {
			continue
		}
		data, err := hex.DecodeString(log.Data)
		if err != nil {
			return nil, err
		}
		argsMap := make(map[string]interface{})
		if err := nEvent.Inputs.UnpackIntoMap(argsMap, data); err != nil {
			return nil, err
		}
		ev := &Event{
			Name:       nEvent.Name,
			RawName:    nEvent.RawName,
			Inputs:     argsMap,
			SMCAddress: log.Address,
			TxHash:     log.TxHash,
		}

		events = append(events, ev)
	}

	return events, nil
}
