package server

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/mongo"
	"github.com/scalarorg/scalar-healer/pkg/openobserve"
	"github.com/scalarorg/scalar-healer/pkg/worker"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Server struct {
	Raw           *echo.Echo
	traceProvider *sdktrace.TracerProvider
	scheduler     worker.Worker
	db            db.DbAdapter
}

func New() *Server {
	config.LoadEnv()

	appName := config.Env.APP_NAME
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
	db := mongo.NewMongoRepository()

	setupAddHandlerEvent(e)
	setupMiddleware(e, db)
	setupErrorHandler(e)
	setupRoute(e)
	setupValidator(e)
	s := setupWorkers()

	return &Server{e, tp, s, db}
}

func (s *Server) Start() error {
	s.printRoutes()

	return s.Raw.Start(fmt.Sprintf("%s:%s", config.Env.API_HOST, config.Env.PORT))
}

func (s *Server) Close() {
	s.db.Close()

	s.scheduler.Shutdown()

	s.Raw.Close()
	err := s.traceProvider.Shutdown(context.Background())
	if err != nil {
		log.Err(err).Msg("Error shutting down trace provider")
	}
}
