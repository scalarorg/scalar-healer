package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/internal/daemon"
	"github.com/scalarorg/scalar-healer/pkg/db/mongo"
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

	db := mongo.NewMongoRepository()
	service := daemon.NewService(config.Env.CLIENTS_CONFIG_PATH, db)
	go func() {
		err := service.Start(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("Cannot start daemon service")
			panic(err)
		}
	}()
}
