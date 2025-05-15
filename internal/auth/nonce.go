package auth

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/pkg/crypto/eip4361"
)

func (h *Handler) GetNonce(c echo.Context) error {
	address := c.Param("address")
	if !common.IsHexAddress(address) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ethereum address")
	}

	msg := eip4361.NewSiweMessage(h.domain, common.HexToAddress(address))
	return c.JSON(http.StatusOK, msg.String())
}
