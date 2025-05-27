package healer

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/utils/slices"
)

func (m *HealerRepository) GetGatewayAddress(ctx context.Context, chainId uint64) (*common.Address, error) {
	address, err := m.Queries.GetGatewayAddress(ctx, db.ConvertUint64ToNumeric(chainId))
	if err != nil {
		return nil, err
	}
	add := common.BytesToAddress(address)
	return &add, nil
}

func (m *HealerRepository) CreateGatewayAddresses(ctx context.Context, addresses [][]byte, chainIds []uint64) error {
	return m.Queries.CreateGatewayAddresses(ctx, sqlc.CreateGatewayAddressesParams{
		Column1: addresses,
		Column2: slices.Map(chainIds, db.ConvertUint64ToNumeric),
	})
}
