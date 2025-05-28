package protocol

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/pkg/db"
)

func ListProtocols(c echo.Context) error {
	db := db.GetRepositoryFromContext(c)
	ctx := c.Request().Context()
	protocols, err := db.GetProtocols(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get protocols")
	}

	return c.JSON(http.StatusOK, protocols)
}
