package healer

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) SaveProtocols(ctx context.Context, protocols []sqlc.Protocol) error {
	var assets []string
	var names []string
	var custodianGroupNames []string
	var custodianGroupUIDs [][]byte
	var tags []string
	var liquidityModels []string
	var symbols []string
	var decimals []int64
	var capacities []pgtype.Numeric
	var dailyMintLimits []pgtype.Numeric
	var avatars []string

	for _, protocol := range protocols {
		assets = append(assets, protocol.Asset)
		names = append(names, protocol.Name)
		custodianGroupNames = append(custodianGroupNames, protocol.CustodianGroupName)
		custodianGroupUIDs = append(custodianGroupUIDs, protocol.CustodianGroupUid)
		tags = append(tags, protocol.Tag)
		liquidityModels = append(liquidityModels, protocol.LiquidityModel)
		symbols = append(symbols, protocol.Symbol)
		decimals = append(decimals, protocol.Decimals)
		capacities = append(capacities, protocol.Capacity)
		dailyMintLimits = append(dailyMintLimits, protocol.DailyMintLimit)
		avatars = append(avatars, protocol.Avatar)
	}

	return m.Queries.SaveProtocols(ctx, sqlc.SaveProtocolsParams{
		Column1:  assets,
		Column2:  names,
		Column3:  custodianGroupNames,
		Column4:  custodianGroupUIDs,
		Column5:  tags,
		Column6:  liquidityModels,
		Column7:  symbols,
		Column8:  decimals,
		Column9:  capacities,
		Column10: dailyMintLimits,
		Column11: avatars,
	})
}

func (m *HealerRepository) GetProtocol(ctx context.Context, name string) (*sqlc.Protocol, error) {
	protocol, err := m.Queries.GetProtocol(ctx, name)
	return &protocol, err
}
