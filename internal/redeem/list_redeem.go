package redeem

import (
	"math"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/internal/middleware"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/utils"
)

type ListRedeemRequest struct {
	Address common.Address `json:"address"`
	Page    int32          `query:"page" validate:"gte=0"`
	Size    int32          `query:"size" validate:"gte=0"`
}

func ListRedeem(c echo.Context) error {
	var body ListRedeemRequest
	body.Address = *middleware.GetAddressFromContext(c)
	if err := utils.BindAndValidate(c, &body); err != nil {
		return err
	}

	db := db.GetRepositoryFromContext(c)
	ctx := c.Request().Context()

	redeemRequests, count, err := db.ListRedeemRequests(ctx, body.Address, body.Page, body.Size)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get redeem requests")
	}

	totalPages := math.Ceil(float64(count) / float64(body.Size))

	return c.JSON(http.StatusOK, echo.Map{
		"items":       redeemRequests,
		"total_pages": totalPages,
	})
}
