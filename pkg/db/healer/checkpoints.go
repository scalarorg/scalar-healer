package healer

import (
	"context"

	"github.com/scalarorg/data-models/scalarnet"
)

func (m *HealerRepository) GetLastEventCheckPoint(ctx context.Context,
	chainName, eventName string, fromBlock uint64) (*scalarnet.EventCheckPoint, error) {
	return nil, nil
}

func (m *HealerRepository) UpdateLastEventCheckPoint(ctx context.Context, lastCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}
