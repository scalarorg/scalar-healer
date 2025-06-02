package healer

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) SaveProtocols(ctx context.Context, protocols []sqlc.Protocol) error {
	var symbols []string
	var names []string
	var custodianGroupNames []string
	var custodianGroupUIDs [][]byte
	var tags []string
	var liquidityModels []string
	var decimals []int64
	var capacities []pgtype.Numeric
	var dailyMintLimits []pgtype.Numeric
	var avatars []string

	for _, protocol := range protocols {
		symbols = append(symbols, protocol.Symbol)
		names = append(names, protocol.Name)
		custodianGroupNames = append(custodianGroupNames, protocol.CustodianGroupName)
		custodianGroupUIDs = append(custodianGroupUIDs, protocol.CustodianGroupUid)
		tags = append(tags, protocol.Tag)
		liquidityModels = append(liquidityModels, protocol.LiquidityModel)
		decimals = append(decimals, protocol.Decimals)
		capacities = append(capacities, protocol.Capacity)
		dailyMintLimits = append(dailyMintLimits, protocol.DailyMintLimit)
		avatars = append(avatars, protocol.Avatar)
	}

	return m.Queries.SaveProtocols(ctx, sqlc.SaveProtocolsParams{
		Column1:  symbols,
		Column2:  names,
		Column3:  custodianGroupNames,
		Column4:  custodianGroupUIDs,
		Column5:  tags,
		Column6:  liquidityModels,
		Column7:  decimals,
		Column8:  capacities,
		Column9:  dailyMintLimits,
		Column10: avatars,
	})
}

func (m *HealerRepository) GetProtocol(ctx context.Context, name string) (*sqlc.Protocol, error) {
	protocol, err := m.Queries.GetProtocol(ctx, name)
	return &protocol, err
}

func (m *HealerRepository) GetProtocols(ctx context.Context) ([]sqlc.ProtocolWithTokenDetails, error) {

	result, err := m.Queries.GetProtocols(ctx)
	mapResult := make(map[string]*sqlc.ProtocolWithTokenDetails)
	for _, r := range result {
		if protocol, ok := mapResult[r.Symbol]; !ok {
			mapResult[r.Symbol] = &sqlc.ProtocolWithTokenDetails{
				Protocol: &sqlc.Protocol{
					ID:                 r.ID,
					Symbol:             r.Symbol,
					Name:               r.Name,
					CustodianGroupName: r.CustodianGroupName,
					CustodianGroupUid:  r.CustodianGroupUid,
					Tag:                r.Tag,
					LiquidityModel:     r.LiquidityModel,
					Decimals:           r.Decimals,
					Avatar:             r.Avatar,
					Capacity:           r.Capacity,
					DailyMintLimit:     r.DailyMintLimit,
					CreatedAt:          r.CreatedAt,
					UpdatedAt:          r.UpdatedAt,
				},
				Tokens: []sqlc.TokenDetails{
					{
						Address: r.Address,
						ChainID: r.ChainID.Int.Int64(),
						Chain:   r.Chain.String,
					},
				},
			}
		} else {
			protocol.Tokens = append(protocol.Tokens, sqlc.TokenDetails{
				Address: r.Address,
				ChainID: r.ChainID.Int.Int64(),
				Chain:   r.Chain.String,
			})
		}
	}

	var protocols []sqlc.ProtocolWithTokenDetails
	for _, v := range mapResult {
		protocols = append(protocols, *v)
	}

	return protocols, err
}
