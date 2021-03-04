// Package kardia
package kardia

import (
	"fmt"

	"go.uber.org/zap"
)

//Nodes like utils struct, which allow user pass trustedNodeUrl and publicNodeUrl
//for better balancing/performance without stress a specific node
type Nodes interface {
}

type nodes struct {
	trusted []Node
	public  []Node

	total  uint64
	logger *zap.Logger
}

type NodesConfig struct {
	TrustedNodeUrls []string
	PublicNodeUrls  []string

	Logger *zap.Logger
}

func NewNodes(cfg NodesConfig) (Nodes, error) {
	nodes := &nodes{}
	for id, url := range cfg.TrustedNodeUrls {
		lgr := cfg.Logger.
			With(zap.String("type", "trusted")).
			With(zap.Int("id", id))
		n, err := NewNode(url, lgr)
		if err != nil {
			zap.S().Warn("cannot connect to url", url)
			continue
		}
		nodes.trusted = append(nodes.trusted, n)
		nodes.total++
	}

	for id, url := range cfg.PublicNodeUrls {
		lgr := cfg.Logger.
			With(zap.String("type", "public")).
			With(zap.Int("id", id))
		n, err := NewNode(url, lgr)
		if err != nil {
			zap.S().Warn("cannot connect to url", url)
			continue
		}
		nodes.trusted = append(nodes.trusted, n)
		nodes.total++
	}

	if nodes.total == 0 {
		return nil, fmt.Errorf("no node available")
	}

	return nodes, nil
}
