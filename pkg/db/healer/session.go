package healer

import (
	"context"

	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/scalar-healer/pkg/db"
)

func (m *HealerRepository) GetRedeemSession(ctx context.Context, chainId string, groupUid string) *db.RedeemSession {
	return nil
}

func (m *HealerRepository) SaveSwitchPhaseValue(ctx context.Context, event *chains.SwitchedPhase) error {
	return nil
}
