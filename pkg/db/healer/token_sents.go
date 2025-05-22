package healer

import (
	"context"

	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/data-models/scalarnet"
)

func (m *HealerRepository) SaveTokenSent(ctx context.Context, tokenSent *chains.TokenSent, eventCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}

func (m *HealerRepository) SaveTokenSents(ctx context.Context, tokenSents []chains.TokenSent) error {
	return nil
}

func (m *HealerRepository) FindPendingBtcTokenSent(ctx context.Context, chainId string, expectedConfirmBlock int32) ([]chains.TokenSent, error) {
	return nil, nil
}
