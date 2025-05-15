package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/pkg/crypto/eip4361"
	"github.com/scalarorg/scalar-healer/pkg/session"
	"github.com/scalarorg/scalar-healer/pkg/utils"
)

type SignInRequest struct {
	Message   string `json:"message" validate:"required"`
	Signature string `json:"signature" validate:"required,hexadecimal"`
}

func (h *Handler) SignIn(c echo.Context) error {
	var body SignInRequest
	if err := utils.BindAndValidate(c, &body); err != nil {
		return err
	}

	msg, err := eip4361.Validate(body.Message, body.Signature)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid signature")
	}

	token, err := session.CreateToken(msg.GetAddress())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create token")
	}

	return c.JSON(http.StatusOK, token)
}
