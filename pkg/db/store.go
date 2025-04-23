package db

import "github.com/jackc/pgx/v5/pgxpool"

var DB *sqlStore

type sqlStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) *sqlStore {
	return &sqlStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
