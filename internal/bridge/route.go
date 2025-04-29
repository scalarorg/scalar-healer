package bridge

import "github.com/labstack/echo/v4"

func Route(g *echo.Group, path string) {
	gr := g.Group(path)
	gr.POST("", CreateBridge)
	gr.GET("/:address", ListBridge)
}
