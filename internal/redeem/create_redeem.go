package redeem

import (
	"encoding/hex"
	"net/http"

	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/internal/middleware"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/utils"
)

type CreateRedeemRequest struct {
	Address       common.Address `json:"address"`
	SourceChain   string         `json:"source_chain" validate:"required"`
	DestChain     string         `json:"dest_chain" validate:"required"`
	Symbol        string         `json:"symbol" validate:"required"`
	Amount        string         `json:"amount" validate:"required"` // bigint format
	LockingScript string         `json:"locking_script" validate:"hexadecimal"`
}

func CreateRedeem(c echo.Context) error {
	var body CreateRedeemRequest
	body.Address = *middleware.GetAddressFromContext(c)
	if err := utils.BindAndValidate(c, &body); err != nil {
		return err
	}

	ctx := c.Request().Context()

	db := db.GetRepositoryFromContext(c)

	_, err := db.GetTokenAddressBySymbol(ctx, body.SourceChain, body.Symbol)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, constants.ErrTokenNotExists)
	}

	// TODO: validate the balance on evm network
	amountz, ok := utils.StringToBigInt(body.Amount)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, constants.ErrInvalidAmount)
	}

	// TODO: validate the locking script
	lockScript, err := hex.DecodeString(body.LockingScript)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, constants.ErrInvalidLockingScript)
	}

	class := txscript.GetScriptClass(lockScript)
	if class == txscript.NonStandardTy {
		return echo.NewHTTPError(http.StatusBadRequest, constants.ErrInvalidLockingScript)
	}

	// txscript.GetWitnessProgramInfo()

	// Save redeem request
	err = db.SaveRedeemRequest(ctx, body.SourceChain, body.DestChain, body.Address, amountz, body.Symbol, lockScript)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save redeem request")
	}

	return c.NoContent(http.StatusOK)
}
