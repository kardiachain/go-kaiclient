// Package kardia
package kardia

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKRC20(t *testing.T) {
	node, err := SetupNodeClient()
	assert.Nil(t, err)
	krc20, err := NewKRC20(node, "0xaf078b4d8694Fe109Ab3F7b2a543C28c5f6EAd16", "0x4f36A53DC32272b97Ae5FF511387E2741D727bdb")
	assert.Nil(t, err)
	fmt.Println("krc20", krc20)

}
