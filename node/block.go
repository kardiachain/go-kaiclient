package node

import (
	"fmt"

	"github.com/kardiachain/go-kaiclient/types"
)

func (e *Eth) getBlock(method string, args ...interface{}) (*types.Block, error) {
	var raw types.Block
	err := e.c.Call(method, &raw, args...)
	if err != nil {
		return nil, err
	}
	fmt.Println("Raw", raw)
	return &raw, nil

}
