package redeem

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/pkg/db/mongo"
	"github.com/scalarorg/scalar-healer/pkg/utils"
)

type ListRedeemRequest struct {
	Address string `param:"address" validate:"eth_addr"`
	Page    int    `query:"page" validate:"gte=0"`
	Size    int    `query:"size" validate:"gte=0"`
}

func ListRedeem(c echo.Context) error {

	var body ListRedeemRequest
	if err := utils.BindAndValidate(c, &body); err != nil {
		return err
	}

	ctx := c.Request().Context()

	db := mongo.GetRepositoryFromContext(c)
	redeemRequests, err := db.ListRedeemRequests(ctx, common.HexToAddress(body.Address), body.Page, body.Size)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get redeem requests")
	}
	return c.JSON(http.StatusOK, redeemRequests)
}
