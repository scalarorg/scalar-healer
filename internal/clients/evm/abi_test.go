package evm_test

import (
	"fmt"
	"testing"

	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/internal/clients/evm"
	"github.com/stretchr/testify/require"
)

func TestGetEventByName(t *testing.T) {
	redeemEvent, ok := evm.GetEventByName(config.EVENT_EVM_REDEEM_TOKEN)
	require.True(t, ok)
	require.NotNil(t, redeemEvent)
	fmt.Printf("redeemEvent %v\n", redeemEvent)
}
