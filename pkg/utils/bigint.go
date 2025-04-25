package utils

import "math/big"

func StringToBigInt(s string) (*big.Int, bool) {
	return big.NewInt(0).SetString(s, 10)
}
