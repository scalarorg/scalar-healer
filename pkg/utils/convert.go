package utils

import (
	"math/big"
	"strconv"
)

func StringToBigInt(s string) (*big.Int, bool) {
	return big.NewInt(0).SetString(s, 10)
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}
