package indexer

import (
	"context"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

type IndexerRepository struct {
	connPool *pgxpool.Pool
	*sqlc.Queries
}

var _ db.IndexerAdapter = (*IndexerRepository)(nil)

type ConnConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

func NewRepository(ctx context.Context, cfg *ConnConfig) *IndexerRepository {
	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	configs, err := pgxpool.ParseConfig(dbSource)

	connPool, err := pgxpool.NewWithConfig(ctx, configs)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	if err = connPool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("cannot ping db")
	}

	return &IndexerRepository{
		connPool: connPool,
		Queries:  sqlc.New(connPool),
	}
}

func (r *IndexerRepository) Close() {
	r.connPool.Close()
}
