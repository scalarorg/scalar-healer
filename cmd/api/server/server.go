package server

import (
	"context"
	"fmt"

	"github.com/0xdavid7/goes-template/config"
	"github.com/0xdavid7/goes-template/pkg/db"
	"github.com/0xdavid7/goes-template/pkg/openobserve"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Server struct {
	Raw           *echo.Echo
	traceProvider *sdktrace.TracerProvider
}

func New() *Server {
	config.LoadEnv()

	appName := config.Env.APP_NAME
	if config.Env.IsDev {
		appName = appName + "-dev"
	}
	openobserve.Init(openobserve.OpenObserveConfig{
		Endpoint:    config.Env.OPENOBSERVE_ENDPOINT,
		Credential:  config.Env.OPENOBSERVE_CREDENTIAL,
		ServiceName: appName,
		Env:         config.Env.ENV,
	})

	config.InitLogger()

	e := echo.New()
	e.HideBanner = true
	tp := openobserve.SetupTraceHTTP()

	setupAddHandlerEvent(e)
	setupMiddleware(e)
	setupErrorHandler(e)
	setupRoute(e)
	setupValidator(e)

	return &Server{e, tp}
}

func (s *Server) Start() error {
	ctx := context.Background()
	loadSvcs(ctx)
	s.printRoutes()

	return s.Raw.Start(fmt.Sprintf("%s:%s", config.Env.API_HOST, config.Env.PORT))
}

func (s *Server) Close() {
	closeSvcs()
	s.Raw.Close()
	err := s.traceProvider.Shutdown(context.Background())
	if err != nil {
		log.Err(err).Msg("Error shutting down trace provider")
	}
}

func loadSvcs(ctx context.Context) {
	db.Init(ctx, config.Env.POSTGRES_URL, config.Env.MIGRATION_URL)
}

func closeSvcs() {
	db.Close()
}
