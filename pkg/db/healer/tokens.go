package healer

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) SaveTokens(ctx context.Context, tokens []sqlc.Token) error {

	var addresses [][]byte
	var chainIds []pgtype.Numeric
	var symbols []string
	var decimals []pgtype.Numeric
	var names []string
	var avatars []string
	var actives []bool

	for _, token := range tokens {
		addresses = append(addresses, token.Address)
		chainIds = append(chainIds, token.ChainID)
		symbols = append(symbols, token.Symbol)
		decimals = append(decimals, token.Decimal)
		names = append(names, token.Name)
		avatars = append(avatars, token.Avatar)
		actives = append(actives, token.Active)
	}

	return m.Queries.SaveTokens(ctx, sqlc.SaveTokensParams{
		Column1: addresses, // address
		Column2: chainIds,  // chain_id,
		Column3: symbols,   // symbol
		Column4: decimals,  // decimal
		Column5: names,     // name
		Column6: avatars,   // avatar
		Column7: actives,   // active,
	})
}

func (m *HealerRepository) GetTokenSymbolByAddress(ctx context.Context, chainId uint64, tokenAddress *common.Address) (string, error) {
	return m.Queries.GetTokenSymbolByAddress(ctx, sqlc.GetTokenSymbolByAddressParams{
		ChainID: db.ConvertUint64ToNumeric(chainId),
		Address: tokenAddress.Bytes(),
	})
}

func (m *HealerRepository) GetTokenAddressBySymbol(ctx context.Context, chainId uint64, tokenSymbol string) (*common.Address, error) {
	address, err := m.Queries.GetTokenAddressBySymbol(ctx, sqlc.GetTokenAddressBySymbolParams{
		ChainID: db.ConvertUint64ToNumeric(chainId),
		Symbol:  tokenSymbol,
	})
	if err != nil {
		tokens, _ := m.Queries.ListTokens(ctx)
		fmt.Printf("tokens: %v\n", tokens)
		return nil, err
	}
	addr := common.BytesToAddress(address)
	return &addr, nil
}
