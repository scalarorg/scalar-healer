package transfer

import "github.com/labstack/echo/v4"

func Route(g *echo.Group, path string) {
	gr := g.Group(path)
	gr.POST("", CreateTransfer)
	gr.GET("/:address", ListTransfer)
}
