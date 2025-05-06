package postgres

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

type PostgresRepository struct {
	connPool *pgxpool.Pool
	*sqlc.Queries
}

var _ db.DbAdapter = (*PostgresRepository)(nil)

type ConnConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

func NewRepository(ctx context.Context, cfg *ConnConfig, migrationURL string) *PostgresRepository {
	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	connPool, err := pgxpool.New(ctx, dbSource)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	if err = connPool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("cannot ping db")
	}

	runDBMigration(migrationURL, dbSource)

	log.Info().Msg("db initialized successfully")

	return &PostgresRepository{
		connPool: connPool,
		Queries:  sqlc.New(connPool),
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func (r *PostgresRepository) Close() {
	r.connPool.Close()
}

func (r *PostgresRepository) DropSchema(name string) {
	r.connPool.Exec(context.Background(), fmt.Sprintf("DROP SCHEMA %s CASCADE", name))
}

func (r *PostgresRepository) ExecQuery(ctx context.Context, query string, args ...interface{}) error {
	_, err := r.connPool.Exec(ctx, query, args...)
	return err
}
