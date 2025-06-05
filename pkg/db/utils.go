package db

import (
	"github.com/labstack/echo/v4"
)

func GetRepositoryContextKey() string {
	return "repository"
}

func GetRepositoryFromContext(c echo.Context) HealderAdapter {
	return c.Get(GetRepositoryContextKey()).(HealderAdapter)
}

func SetRepositoryToContext(c echo.Context, next echo.HandlerFunc, m HealderAdapter) error {
	c.Set(GetRepositoryContextKey(), m)
	return next(c)
}
