package middleware

import (
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/pkg/session"
)

func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
		}

		address, err := session.ValidateToken(parts[1])
		if err != nil {
			log.Err(err).Msg("invalid token")
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		c.Set(constants.AUTH_ADDRESS_KEY, address)

		return next(c)
	}
}

func GetAddress(c echo.Context) *common.Address {
	address := c.Get(constants.AUTH_ADDRESS_KEY)
	if address == nil {
		return nil
	}
	return address.(*common.Address)
}
