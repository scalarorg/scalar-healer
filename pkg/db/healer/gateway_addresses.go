package healer

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) GetGatewayAddress(ctx context.Context, chain string) (*common.Address, error) {
	address, err := m.Queries.GetGatewayAddress(ctx, chain)
	if err != nil {
		return nil, err
	}
	add := common.BytesToAddress(address)
	return &add, nil
}

func (m *HealerRepository) CreateGatewayAddresses(ctx context.Context, addresses [][]byte, chains []string) error {
	return m.Queries.CreateGatewayAddresses(ctx, sqlc.CreateGatewayAddressesParams{
		Column1: addresses,
		Column2: chains,
	})
}
