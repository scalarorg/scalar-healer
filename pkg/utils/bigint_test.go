package utils_test

import (
	"testing"

	"github.com/scalarorg/scalar-healer/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestStringToBigInt(t *testing.T) {
	result, ok := utils.StringToBigInt("1aaaa")
	require.False(t, ok)
	require.Nil(t, result)
}
