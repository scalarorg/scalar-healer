package healer

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

type contextKey string

const txContextKey contextKey = "transaction_context"

func (store *HealerRepository) requireTx(ctx context.Context, fn func(ctx context.Context) error) error {

	log.Info().Msgf("requireTx: %v", ctx.Value(txContextKey))
	
	if ctx.Value(txContextKey) == nil {
		return fmt.Errorf("must be called within execTx transaction")
	}


	if !ctx.Value(txContextKey).(bool) {
		return fmt.Errorf("must be called within execTx transaction")
	}

	return fn(ctx)
}

func (store *HealerRepository) execTx(ctx context.Context, fn func(context.Context, *sqlc.Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)

	ctx = context.WithValue(ctx, txContextKey, true)

	err = fn(ctx, q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
