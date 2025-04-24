package server

import (
	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/internal/health"
)

func setupRoute(e *echo.Echo) {
	health.Route(e, "/health")
	api := e.Group("/api")
	_ = api
}
