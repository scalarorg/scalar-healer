package postgres

import (
	"context"

	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/data-models/scalarnet"
)

func (m *PostgresRepository) SaveTokenSent(ctx context.Context, tokenSent *chains.TokenSent, eventCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}

func (m *PostgresRepository) SaveTokenSents(ctx context.Context, tokenSents []chains.TokenSent) error {
	return nil
}

func (m *PostgresRepository) FindPendingBtcTokenSent(ctx context.Context, chainId string, expectedConfirmBlock int32) ([]chains.TokenSent, error) {
	return nil, nil
}
