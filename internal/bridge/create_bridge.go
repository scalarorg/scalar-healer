package bridge

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/pkg/crypto/eip712"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/utils"
)

type CreateBridgeRequest struct {
	eip712.BaseRequest
	TxHash string `json:"tx_hash" validate:"hexadecimal,len=66"`
}

func CreateBridge(c echo.Context) error {
	var body CreateBridgeRequest
	if err := utils.BindAndValidate(c, &body); err != nil {
		return err
	}

	ctx := c.Request().Context()

	db := db.GetRepositoryFromContext(c)

	gatewayAddress, err := db.GetGatewayAddress(ctx, body.ChainID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, constants.ErrNotFoundGateway)

	}

	txHash := common.HexToHash(body.TxHash)

	message := eip712.NewBridgeRequestMessage(&body.BaseRequest, txHash)
	err = message.Validate(ctx, gatewayAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Save redeem request
	err = db.SaveBridgeRequest(ctx, body.ChainID, common.HexToAddress(body.Address), common.FromHex(body.Signature), txHash.Bytes(), body.Nonce)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save bridge request")
	}

	return c.NoContent(http.StatusOK)
}
