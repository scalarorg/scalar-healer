package redeem

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetRedeemRequests(c echo.Context) error {

	return c.NoContent(http.StatusOK)
}
