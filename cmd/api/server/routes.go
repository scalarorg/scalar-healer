package server

import (
	"github.com/0xdavid7/goes-template/internal/health"
	"github.com/labstack/echo/v4"
)

func setupRoute(e *echo.Echo) {
	health.Route(e, "/health")
	api := e.Group("/api")
	_ = api
}
