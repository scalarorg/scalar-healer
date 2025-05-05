package db

import (
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func GetRepositoryContextKey() string {
	return "repository"
}

func GetRepositoryFromContext(c echo.Context) DbAdapter {
	return c.Get(GetRepositoryContextKey()).(DbAdapter)
}

func SetRepositoryToContext(c echo.Context, next echo.HandlerFunc, m DbAdapter) error {
	c.Set(GetRepositoryContextKey(), m)
	return next(c)
}

func ConvertUint64ToNumeric(n uint64) pgtype.Numeric {
	var binary pgtype.Numeric
	var big = &big.Int{}
	big.SetUint64(n)
	binary.Scan(big.String())
	return binary
}
