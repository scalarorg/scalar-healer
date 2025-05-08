package testutils

import (
	"context"
	"time"

	"log"

	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/postgres"
	"github.com/scalarorg/scalar-healer/pkg/utils"
	"github.com/testcontainers/testcontainers-go"
	pgc "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RunWithTestDBFunc func(ctx context.Context, repo db.DbAdapter) error

func RunWithTestDB(callback RunWithTestDBFunc) {
	ctx := context.Background()
	config.LoadEnvWithPath(".env")

	rootPath, err := utils.RootPath()
	if err != nil {
		log.Fatal("Failed to get root path")
	}

	config.Env.MIGRATION_URL = "file://" + rootPath + "/pkg/db/sqlc/migration"
	config.Env.POSTGRES_DB = "scalar_healer_test"
	config.Env.POSTGRES_USER = "postgres"
	config.Env.POSTGRES_PASSWORD = "postgres"
	config.Env.ENV = "test"

	pgContainer, err := pgc.Run(ctx,
		"postgres:latest",
		pgc.WithDatabase(config.Env.POSTGRES_DB),
		pgc.WithUsername(config.Env.POSTGRES_USER),
		pgc.WithPassword(config.Env.POSTGRES_PASSWORD),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatal("Failed to start postgres container")
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatal("Failed to get postgres port")
	}

	var repo *postgres.PostgresRepository

	for retries := 0; retries < 3; retries++ {
		repo = postgres.NewRepository(ctx, &postgres.ConnConfig{
			User:     config.Env.POSTGRES_USER,
			Password: config.Env.POSTGRES_PASSWORD,
			Host:     "127.0.0.1", // Use IP instead of localhost
			Port:     port.Int(),
			DBName:   config.Env.POSTGRES_DB,
		}, config.Env.MIGRATION_URL)

		if repo != nil {
			break
		}

		log.Println("retrying database connection")
		time.Sleep(time.Second * 2)
	}

	if repo == nil {
		log.Fatal("failed to connect to database after retries")
	}

	err = callback(ctx, repo)
	if err != nil {
		log.Fatal("Error in callback")
	}

	if repo != nil {
		repo.Close()
	}

	if pgContainer != nil {
		log.Println("Terminating postgres container")
		err = pgContainer.Terminate(ctx)
		if err != nil {
			log.Println("Failed to terminate postgres container")
		}
	}
}
