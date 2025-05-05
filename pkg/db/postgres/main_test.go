package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db/postgres"
)

var (
	repo *postgres.PostgresRepository
)

func TestMain(m *testing.M) {
	config.LoadEnvWithPath("../../../.env")

	repo = postgres.NewRepository(context.Background(), &postgres.ConnConfig{
		User:     config.Env.POSTGRES_USER,
		Password: config.Env.POSTGRES_PASSWORD,
		Host:     config.Env.POSTGRES_HOST,
		Port:     config.Env.POSTGRES_PORT,
		DBName:   config.Env.POSTGRES_DB,
	}, config.Env.MIGRATION_URL)
	code := m.Run()
	cleanupTestDB()
	os.Exit(code)
}

func cleanupTestDB() {
	repo.Close()
}
