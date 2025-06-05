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
	return num.Int.Uint64()
}
