package postgres

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

func (m *PostgresRepository) GetGatewayAddress(ctx context.Context, chainId uint64) (*common.Address, error) {
	// var data struct {
	// 	Address []byte `bson:"address"`
	// }

	// opts := options.FindOne().SetProjection(bson.M{
	// 	"address": 1,
	// })
	// err := m.GatewayAddresses.FindOne(ctx, map[string]interface{}{
	// 	"chain_id": chainId,
	// }, opts).Decode(&data)
	// if err != nil {
	// 	return nil, err
	// }

	// add := common.BytesToAddress(data.Address)
	// return &add, nil
	return nil, nil
}
