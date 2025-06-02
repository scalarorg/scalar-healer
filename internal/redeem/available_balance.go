package redeem

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetAvailableBalanceRequest struct {
	LockingScript string `param:"locking_script" validate:"hexadecimal"`
	Symbol        string `param:"symbol" validate:"required"`
}

// TODO: use FFI to bitcoin-vault to derive the locking address and then get available balance
func GetAvailableBalance(c echo.Context) error {
	// var body GetAvailableBalanceRequest
	// if err := utils.BindAndValidate(c, &body); err != nil {
	// 	return err
	// }

	// db := db.GetRepositoryFromContext(c)
	// ctx := c.Request().Context()

	// protocol, err := db.GetProtocol(ctx, body.Symbol)
	// if err != nil {
	// 	return err
	// }

	// err = utils.ValidateLockingScript(body.LockingScript)

	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	// }

	// upcLockingScript, err := hex.DecodeString(body.LockingScript)

	return c.JSON(http.StatusOK, "100")
}
