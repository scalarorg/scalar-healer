package mongo

import (
	"context"

	"github.com/scalarorg/data-models/scalarnet"
)

func (m *MongoRepository) GetLastEventCheckPoint(ctx context.Context,
	chainName, eventName string, fromBlock uint64) (*scalarnet.EventCheckPoint, error) {
	return nil, nil
}

func (m *MongoRepository) UpdateLastEventCheckPoint(ctx context.Context, lastCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}
