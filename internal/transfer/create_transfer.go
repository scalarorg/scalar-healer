package transfer

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/pkg/crypto/eip712"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/utils"
)

type CreateTransferRequest struct {
	eip712.BaseRequest
	DestinationChain   string `json:"destination_chain" validate:"required"`
	DestinationAddress string `json:"destination_address" validate:"eth_addr"`
	Symbol             string `json:"symbol" validate:"required"`
	Amount             string `json:"amount" validate:"required"`
}

func CreateTransfer(c echo.Context) error {
	var body CreateTransferRequest
	if err := utils.BindAndValidate(c, &body); err != nil {
		return err
	}

	ctx := c.Request().Context()

	db := db.GetRepositoryFromContext(c)

	gatewayAddress, err := db.GetGatewayAddress(ctx, body.Chain)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, constants.ErrNotFoundGateway)

	}

	// TODO: validate body.symbol is valid on destination chain as well
	_, err = db.GetTokenAddressBySymbol(ctx, body.Chain, body.Symbol)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, constants.ErrTokenNotExists)
	}

	destAddress := common.HexToAddress(body.DestinationAddress)
	amountz, ok := utils.StringToBigInt(body.Amount)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, constants.ErrInvalidAmount)
	}

	message := eip712.NewTransferRequestMessage(&body.BaseRequest, body.DestinationChain, &destAddress, body.Symbol, amountz)
	err = message.Validate(ctx, gatewayAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Save redeem request
	err = db.SaveTransferRequest(ctx, body.Chain, body.Address, common.FromHex(body.Signature), amountz, body.DestinationChain, &destAddress, body.Symbol, body.Nonce)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save bridge request")
	}

	return c.NoContent(http.StatusOK)
}
