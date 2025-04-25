package redeem

import "github.com/labstack/echo/v4"

func Route(g *echo.Group, path string) {
	redeemGr := g.Group(path)
	redeemGr.POST("", CreateRedeem)
	redeemGr.GET("/:address", GetRedeemRequests)
}
