package sqlc

import (
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertUint64ToNumeric(n uint64) pgtype.Numeric {
	var binary pgtype.Numeric
	var big = &big.Int{}
	big.SetUint64(n)
	binary.Scan(big.String())
	return binary
}

func ConvertNumericToUint64(num pgtype.Numeric) uint64 {
	if !num.Valid || num.NaN || num.InfinityModifier != 0 || num.Int == nil {
		return 0
	}
	val := new(big.Int).Set(num.Int)
	if num.Exp < 0 {
		exp := big.NewInt(1)
		exp.Exp(big.NewInt(10), big.NewInt(int64(-num.Exp)), nil)
		val.Div(val, exp)
	} else if num.Exp > 0 {
		exp := big.NewInt(1)
		exp.Exp(big.NewInt(10), big.NewInt(int64(num.Exp)), nil)
		val.Mul(val, exp)
	}
	if val.Sign() < 0 {
		return 0
	}
	return val.Uint64()
}