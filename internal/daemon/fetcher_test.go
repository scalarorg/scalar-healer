package daemon_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/scalarorg/scalar-healer/internal/daemon"
	"github.com/stretchr/testify/assert"
)

func TestGetLastestBlock(t *testing.T) {
	block, err := daemon.GetCurrentBlock()
	assert.NoError(t, err)

	fmt.Println(block)
}

func TestUtxoList(t *testing.T) {
	hex, err := hex.DecodeString("5120a8fc50d87f16d892b4d4d087d259c0ab417e106b044b291a7728d2ae1343de7f")
	assert.NoError(t, err)

	utxoList, _, err := daemon.GetUtxoList(hex, []byte{1, 2, 3})
	assert.NoError(t, err)

	fmt.Println(len(utxoList))
}
