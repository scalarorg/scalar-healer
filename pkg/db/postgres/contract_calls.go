package postgres

import (
	"context"

	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/data-models/scalarnet"
)

func (m *PostgresRepository) SaveContractCallWithToken(ctx context.Context, contractCallWithToken *chains.ContractCallWithToken, eventCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}
