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
