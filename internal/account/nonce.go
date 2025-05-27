package account

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/internal/middleware"
	"github.com/scalarorg/scalar-healer/pkg/db"
)

func GetNonce(c echo.Context) error {

	address := middleware.GetAddressFromContext(c)

	db := db.GetRepositoryFromContext(c)
	ctx := c.Request().Context()

	return c.JSON(http.StatusOK, db.GetNonce(ctx, *address))
}
