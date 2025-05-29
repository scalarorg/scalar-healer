package redeem

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/internal/middleware"
	"github.com/scalarorg/scalar-healer/pkg/crypto/eip712"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/utils"
)

type CreateRedeemRequest struct {
	eip712.BaseRequest
	Symbol  string `json:"symbol" validate:"required"`
	Amount  string `json:"amount" validate:"required"` // bigint format
}

func CreateRedeem(c echo.Context) error {
	var body CreateRedeemRequest
	body.Address = *middleware.GetAddressFromContext(c)
	if err := utils.BindAndValidate(c, &body); err != nil {
		return err
	}

	ctx := c.Request().Context()

	db := db.GetRepositoryFromContext(c)

	gatewayAddress, err := db.GetGatewayAddress(ctx, body.ChainID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, constants.ErrNotFoundGateway)
	}

	_, err = db.GetTokenAddressBySymbol(ctx, body.ChainID, body.Symbol)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, constants.ErrTokenNotExists)
	}

	// TODO: validate the balance on evm network
	amountz, ok := utils.StringToBigInt(body.Amount)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, constants.ErrInvalidAmount)
	}

	message := eip712.NewRedeemRequestMessage(&body.BaseRequest, body.Symbol, amountz)
	err = message.Validate(ctx, gatewayAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Save redeem request
	err = db.SaveRedeemRequest(ctx, body.ChainID, body.Address, common.FromHex(body.Signature), amountz, body.Symbol, body.Nonce)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save redeem request")
	}

	return c.NoContent(http.StatusOK)
}
