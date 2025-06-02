package redeem

import (
	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/internal/middleware"
)

func Route(g *echo.Group, path string) {
	gr := g.Group(path, middleware.Authenticate)
	gr.POST("", CreateRedeem)
	gr.POST("/available_balance", GetAvailableBalance)
	gr.GET("/:address", ListRedeem)
}
