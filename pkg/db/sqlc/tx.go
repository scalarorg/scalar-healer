package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TxFunc func(ctx context.Context, q *Queries, args ...interface{}) error

func RequireTx(fn TxFunc) TxFunc {
	return func(ctx context.Context, q *Queries, args ...interface{}) error {
		if _, ok := q.db.(pgx.Tx); !ok {
			panic("function must be called within a transaction (execTx)")
		}
		return fn(ctx, q, args...)
	}
}
