// Package types
package types

import (
	"encoding/json"
)

type CacheBlock struct {
	Hash     string
	Number   uint64
	IsSynced bool
}

func (c CacheBlock) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(data)
}
