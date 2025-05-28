package server

import (
	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/internal/account"
	"github.com/scalarorg/scalar-healer/internal/auth"
	"github.com/scalarorg/scalar-healer/internal/bridge"
	"github.com/scalarorg/scalar-healer/internal/health"
	"github.com/scalarorg/scalar-healer/internal/protocol"
	"github.com/scalarorg/scalar-healer/internal/redeem"
	"github.com/scalarorg/scalar-healer/internal/transfer"
)

func setupRoute(e *echo.Echo) {
	health.Route(e, "/health")
	api := e.Group("/api")
	auth.Route(api, "/auth")
	redeem.Route(api, "/redeem")
	bridge.Route(api, "/bridge")
	transfer.Route(api, "/transfer")
	account.Route(api, "/account")
	protocol.Route(api, "/protocol")
}
