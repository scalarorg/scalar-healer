package healer

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type Tracer struct{}

func (t *Tracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	log.Debug().Str("sql", data.SQL).Interface("args", data.Args).Msg("query start")
	return ctx
}

func (t *Tracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}
