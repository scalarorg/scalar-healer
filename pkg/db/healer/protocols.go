package healer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) SaveProtocols(ctx context.Context, protocols []sqlc.Protocol) error {
	var symbols []string
	var names []string
	var bitconPubkeys [][]byte
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
		bitconPubkeys = append(bitconPubkeys, protocol.BitcoinPubkey)
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
		Column3:  bitconPubkeys,
		Column4:  custodianGroupNames,
		Column5:  custodianGroupUIDs,
		Column6:  tags,
		Column7:  liquidityModels,
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

func (m *HealerRepository) GetProtocols(ctx context.Context) ([]sqlc.ProtocolWithTokenDetails, error) {
	result, err := m.Queries.GetProtocols(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get protocols: %w", err)
	}

	protocolsBySymbol := make(map[string]*sqlc.ProtocolWithTokenDetails)

	for _, row := range result {
		var custodians []sqlc.Custodian
		err := json.Unmarshal(row.Custodians, &custodians)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal custodians: %w", err)
		}

		// Create token details for the current row
		tokenDetails := sqlc.TokenDetails{
			Address: row.Address,
			ChainID: row.ChainID.Int.Int64(),
			Chain:   row.Chain.String,
		}

		// Get or create protocol entry
		protocol, exists := protocolsBySymbol[row.Symbol]
		if !exists {
			// Initialize new protocol with base details
			protocol = &sqlc.ProtocolWithTokenDetails{
				Protocol: &sqlc.Protocol{
					ID:                 row.ID,
					Symbol:             row.Symbol,
					Name:               row.Name,
					BitcoinPubkey:      row.BitcoinPubkey,
					CustodianGroupName: row.CustodianGroupName,
					CustodianGroupUid:  row.CustodianGroupUid,
					Tag:                row.Tag,
					LiquidityModel:     row.LiquidityModel,
					Decimals:           row.Decimals,
					Avatar:             row.Avatar,
					Capacity:           row.Capacity,
					DailyMintLimit:     row.DailyMintLimit,
					CreatedAt:          row.CreatedAt,
					UpdatedAt:          row.UpdatedAt,
				},
				Custodians:      custodians,
				CustodianQuorum: row.Quorum,
				Tokens:          []sqlc.TokenDetails{tokenDetails},
			}
			protocolsBySymbol[row.Symbol] = protocol
		} else {
			// Append token details to existing protocol
			protocol.Tokens = append(protocol.Tokens, tokenDetails)
		}
	}

	// Convert map to slice
	protocols := make([]sqlc.ProtocolWithTokenDetails, 0, len(protocolsBySymbol))
	for _, protocol := range protocolsBySymbol {
		protocols = append(protocols, *protocol)
	}

	return protocols, nil
}
