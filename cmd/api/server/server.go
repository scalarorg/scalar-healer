package server

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/internal/worker"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/openobserve"
	"github.com/scalarorg/scalar-healer/pkg/session"
	"github.com/scalarorg/scalar-healer/pkg/tofnd"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Server struct {
	Raw           *echo.Echo
	DB            db.HealderAdapter
	traceProvider *sdktrace.TracerProvider
	scheduler     *worker.Scheduler
}

func New(db db.HealderAdapter) *Server {
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
	tofndManager := tofnd.NewManager(config.Env.CLIENTS_CONFIG_PATH)

	setupAddHandlerEvent(e)
	setupMiddleware(e, db)
	setupErrorHandler(e)
	setupRoute(e)
	setupValidator(e)
	s := setupWorkers(db, tofndManager)

	return &Server{e, db, tp, s}
}

func (s *Server) Start() error {
	s.printRoutes()

	s.setupSvs()

	s.scheduler.Start()

	return s.Raw.Start(fmt.Sprintf("%s:%s", config.Env.API_HOST, config.Env.PORT))
}

func (s *Server) setupSvs() {
	session.Init([]byte(config.Env.JWT_SECRET), config.Env.JWT_DURATION)
}

func (s *Server) Close() {
	s.DB.Close()

	s.scheduler.Shutdown()

	s.Raw.Close()
	err := s.traceProvider.Shutdown(context.Background())
	if err != nil {
		log.Err(err).Msg("Error shutting down trace provider")
	}
}
