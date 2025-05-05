package postgres

import (
	"context"

	"github.com/scalarorg/data-models/scalarnet"
)

func (m *PostgresRepository) GetLastEventCheckPoint(ctx context.Context,
	chainName, eventName string, fromBlock uint64) (*scalarnet.EventCheckPoint, error) {
	return nil, nil
}

func (m *PostgresRepository) UpdateLastEventCheckPoint(ctx context.Context, lastCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}
