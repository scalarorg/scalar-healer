package mongo

import (
	"context"

	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/data-models/scalarnet"
)

func (m *MongoRepository) SaveTokenSent(ctx context.Context, tokenSent *chains.TokenSent, eventCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}

func (m *MongoRepository) SaveTokenSents(ctx context.Context, tokenSents []*chains.TokenSent) error {
	return nil
}

func (m *MongoRepository) FindPendingBtcTokenSent(ctx context.Context, chainId string, expectedConfirmBlock int) ([]*chains.TokenSent, error) {
	return nil, nil
}
