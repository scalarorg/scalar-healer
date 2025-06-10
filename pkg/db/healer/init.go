package healer

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

type HealerRepository struct {
	connPool *pgxpool.Pool
	*sqlc.Queries
}

var _ db.HealderAdapter = (*HealerRepository)(nil)

type ConnConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

func NewRepository(ctx context.Context, cfg *ConnConfig, migrationURL string) *HealerRepository {
	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	configs, err := pgxpool.ParseConfig(dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot parse db config")
	}
	if config.Env.IS_TEST  {
		configs.ConnConfig.Tracer = &Tracer{}
	}

	connPool, err := pgxpool.NewWithConfig(ctx, configs)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	if err = connPool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("cannot ping db")
	}

	runDBMigration(migrationURL, dbSource)

	log.Info().Msg("db initialized successfully")

	return &HealerRepository{
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

func (r *HealerRepository) Close() {
	r.connPool.Close()
}

func (r *HealerRepository) DropSchema(name string) {
	r.connPool.Exec(context.Background(), fmt.Sprintf("DROP SCHEMA %s CASCADE", name))
}

func (r *HealerRepository) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := r.connPool.Exec(ctx, query, args...)
	return err
}

func (r *HealerRepository) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return r.connPool.Query(ctx, query, args...)
}

func (r *HealerRepository) TruncateTable(ctx context.Context, tableName string) error {
	_, err := r.connPool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s", tableName))
	return err
}

func (r *HealerRepository) TruncateTables(ctx context.Context, tableNames ...string) error {
	if len(tableNames) == 0 {
		return nil
	}
	query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", strings.Join(tableNames, ", "))
	_, err := r.connPool.Exec(ctx, query)
	return err
}
