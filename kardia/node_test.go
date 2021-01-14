// Package kardia
package kardia

import (
	"testing"

	"go.uber.org/zap"
)

var url = "https://dev-1.kardiachain.io"

func SetupNodeClient() (Node, error) {
	lgr, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	node, err := NewNode(url, lgr)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func TestNode_Ping(t *testing.T) {

}

func TestNode_GetBalance(t *testing.T) {

}
