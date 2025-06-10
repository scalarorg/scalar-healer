package healer

import (
	"context"
	"fmt"

	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

type contextKey string

const txContextKey contextKey = "transaction_context"

func (store *HealerRepository) requireTx(ctx context.Context, fn func(ctx context.Context) error) error {

	if ctx.Value(txContextKey) == nil {
		return fmt.Errorf("must be called within execTx transaction")
	}

	if !ctx.Value(txContextKey).(bool) {
		return fmt.Errorf("must be called within execTx transaction")
	}

	return fn(ctx)
}

func (store *HealerRepository) execTx(ctx context.Context, fn func(context.Context, *sqlc.Queries) error) error {
	txCtx := context.WithValue(ctx, txContextKey, true)

	tx, err := store.connPool.Begin(txCtx)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)

	err = fn(txCtx, q)
	if err != nil {
		if rbErr := tx.Rollback(txCtx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(txCtx)
}
