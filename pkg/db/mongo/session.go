package mongo

import (
	"context"

	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/scalar-healer/pkg/db/models"
)

func (m *MongoRepository) GetRedeemSession(ctx context.Context, chainId string, groupUid string) *models.RedeemSession {
	return nil
}

func (m *MongoRepository) SaveSwitchPhaseValue(ctx context.Context, event *chains.SwitchedPhase) error {
	return nil
}
