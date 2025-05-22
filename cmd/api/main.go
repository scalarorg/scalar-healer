package main

import (
	"context"

	"github.com/scalarorg/scalar-healer/cmd/api/server"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db/healer"
)

func main() {
	config.LoadEnv()

	db := healer.NewRepository(context.Background(), &healer.ConnConfig{
		User:     config.Env.HEALER_POSTGRES_USER,
		Password: config.Env.HEALER_POSTGRES_PASSWORD,
		Host:     config.Env.HEALER_POSTGRES_HOST,
		Port:     config.Env.HEALER_POSTGRES_PORT,
		DBName:   config.Env.HEALER_POSTGRES_DB,
	}, config.Env.MIGRATION_URL)

	s := server.New(db)
	defer s.Close()

	err := s.Start()
	if err != nil {
		panic(err)
	}
}
