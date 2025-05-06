package main

import (
	"context"

	"github.com/scalarorg/scalar-healer/cmd/api/server"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db/postgres"
)

func main() {
	config.LoadEnv()

	db := postgres.NewRepository(context.Background(), &postgres.ConnConfig{
		User:     config.Env.POSTGRES_USER,
		Password: config.Env.POSTGRES_PASSWORD,
		Host:     config.Env.POSTGRES_HOST,
		Port:     config.Env.POSTGRES_PORT,
		DBName:   config.Env.POSTGRES_DB,
	}, config.Env.MIGRATION_URL)

	s := server.New(db)
	defer s.Close()

	err := s.Start()
	if err != nil {
		panic(err)
	}
}
