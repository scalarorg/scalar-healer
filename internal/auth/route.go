package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/config"
)

func Route(g *echo.Group, path string) {
	handler := NewHandler(config.Env.AUTH_DOMAIN)
	gr := g.Group(path)
	gr.GET("/nonce/:address", handler.GetNonce)
	gr.POST("/signin", handler.SignIn)
}
