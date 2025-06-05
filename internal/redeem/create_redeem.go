package redeem

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/internal/middleware"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/utils"
	"github.com/scalarorg/scalar-healer/pkg/utils/chains"
)

type CreateRedeemRequest struct {
	Address       common.Address   `json:"address"`
	SourceChain   chains.ChainName `json:"source_chain" validate:"required"`
	DestChain     chains.ChainName `json:"dest_chain" validate:"required"`
	Symbol        string           `json:"symbol" validate:"required"`
	Amount        string           `json:"amount" validate:"required"` // bigint format
	LockingScript string           `json:"locking_script" validate:"hexadecimal"`
}

func CreateRedeem(c echo.Context) error {
	var body CreateRedeemRequest
	body.Address = *middleware.GetAddressFromContext(c)
	if err := utils.BindAndValidate(c, &body); err != nil {
		return err
	}

	ctx := c.Request().Context()

	db := db.GetRepositoryFromContext(c)

	// TODO: validate the balance on evm network
	amountz, ok := utils.StringToBigInt(body.Amount)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, constants.ErrInvalidAmount)
	}

	lockScript, err := utils.ValidateLockingScript(body.LockingScript)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, constants.ErrInvalidLockingScript)
	}

	if !body.SourceChain.IsEvmChain() {
		return echo.NewHTTPError(http.StatusBadRequest, "source chain is not a EVM chain")
	}

	if !body.DestChain.IsBitcoinChain() {
		return echo.NewHTTPError(http.StatusBadRequest, "destination chain is not a bitcoin chain")
	}

	// Save redeem request
	err = db.SaveRedeemRequest(ctx, body.SourceChain.String(), body.DestChain.String(), body.Address, amountz, body.Symbol, lockScript)
	if err != nil {
		log.Error().Err(err).Msg("failed to save redeem request")
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save redeem request")
	}

	return c.NoContent(http.StatusOK)
}
