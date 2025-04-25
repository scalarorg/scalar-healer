package utils_test

import (
	"testing"

	"github.com/scalarorg/scalar-healer/pkg/utils"
)

func TestStringToBigInt(t *testing.T) {
	result, ok := utils.StringToBigInt("1aaaa")
	if !ok {
		t.Errorf("Expected ok to be true")
	}

	t.Log(result)
}
