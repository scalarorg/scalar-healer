package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/internal/daemon"
	"github.com/scalarorg/scalar-healer/pkg/db/combined"
	"github.com/scalarorg/scalar-healer/pkg/db/healer"
	"github.com/scalarorg/scalar-healer/pkg/db/indexer"
	"github.com/scalarorg/scalar-healer/pkg/openobserve"
)

func main() {
	config.LoadEnv()

	appName := config.Env.APP_NAME
	openobserve.Init(openobserve.OpenObserveConfig{
		Endpoint:    config.Env.OPENOBSERVE_ENDPOINT,
		Credential:  config.Env.OPENOBSERVE_CREDENTIAL,
		ServiceName: appName,
		Env:         config.Env.ENV,
	})

	config.InitLogger()

	healerDb := healer.NewRepository(context.Background(), &healer.ConnConfig{
		User:     config.Env.HEALER_POSTGRES_USER,
		Password: config.Env.HEALER_POSTGRES_PASSWORD,
		Host:     config.Env.HEALER_POSTGRES_HOST,
		Port:     config.Env.HEALER_POSTGRES_PORT,
		DBName:   config.Env.HEALER_POSTGRES_DB,
	}, config.Env.MIGRATION_URL)

	indexerDb := indexer.NewRepository(context.Background(), &indexer.ConnConfig{
		User:     config.Env.INDEXER_POSTGRES_USER,
		Password: config.Env.INDEXER_POSTGRES_PASSWORD,
		Host:     config.Env.INDEXER_POSTGRES_HOST,
		Port:     config.Env.INDEXER_POSTGRES_PORT,
		DBName:   config.Env.INDEXER_POSTGRES_DB,
	})

	db := combined.NewCombinedManager(healerDb, indexerDb)

	service := daemon.NewService(config.Env.CLIENTS_CONFIG_PATH, config.Env.EVM_PRIVATE_KEY, db)
	err := service.Start(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Cannot start daemon service")
		panic(err)
	}
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down relayer...")
	service.Stop()
}
