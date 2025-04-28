package mongo

import (
	"context"

	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/data-models/scalarnet"
)

func (m *MongoRepository) SaveContractCallWithToken(ctx context.Context, contractCallWithToken *chains.ContractCallWithToken, eventCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}
