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
	krc20, err := NewKRC20(node, "0x443c6B6240AFDf1D7497aAB64858F5b1363d321B", "0x4f36A53DC32272b97Ae5FF511387E2741D727bdb")
	assert.Nil(t, err)
	fmt.Println("krc20", krc20)

}
