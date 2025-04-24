package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/cmd/api/server"
	"github.com/scalarorg/scalar-healer/internal/daemon"
)

func main() {
	s := server.New()
	service := daemon.NewService()
	go func() {
		err := service.Start(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("Cannot start daemon service")
			panic(err)
		}
	}()
	err := s.Start()
	defer s.Close()
	if err != nil {
		panic(err)
	}
}
