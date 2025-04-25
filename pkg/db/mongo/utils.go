package mongo

import (
	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/pkg/db"
)

func GetRepositoryContextKey() string {
	return "mongo_repository"
}

func GetRepositoryFromContext(c echo.Context) db.DbAdapter {
	return c.Get(GetRepositoryContextKey()).(*MongoRepository)
}

func SetRepositoryToContext(c echo.Context, next echo.HandlerFunc, m db.DbAdapter) error {
	c.Set(GetRepositoryContextKey(), m)
	return next(c)
}
