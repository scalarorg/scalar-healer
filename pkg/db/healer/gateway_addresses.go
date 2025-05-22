package healer

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/pkg/db"
)

func (m *HealerRepository) GetGatewayAddress(ctx context.Context, chainId uint64) (*common.Address, error) {
	address, err := m.Queries.GetGatewayAddress(ctx, db.ConvertUint64ToNumeric(chainId))
	if err != nil {
		return nil, err
	}
	add := common.BytesToAddress(address)
	return &add, nil
}
