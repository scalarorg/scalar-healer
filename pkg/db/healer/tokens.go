package healer

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) SaveTokens(ctx context.Context, tokens []sqlc.Token) error {

	var addresses [][]byte
	var chains []string
	var chainIds []pgtype.Numeric
	var symbols []string
	var actives []bool

	for _, token := range tokens {
		addresses = append(addresses, token.Address)
		chains = append(chains, token.Chain)
		chainIds = append(chainIds, token.ChainID)
		symbols = append(symbols, token.Symbol)
		actives = append(actives, token.Active)
	}

	return m.Queries.SaveTokens(ctx, sqlc.SaveTokensParams{
		Column1: addresses, // address
		Column2: chains,    // chain,
		Column3: chainIds,  // chain_id,
		Column4: symbols,   // symbol
		Column5: actives,   // active,
	})
}

func (m *HealerRepository) GetTokenSymbolByAddress(ctx context.Context, chain string, tokenAddress *common.Address) (string, error) {
	return m.Queries.GetTokenSymbolByAddress(ctx, sqlc.GetTokenSymbolByAddressParams{
		Chain:   chain,
		Address: tokenAddress.Bytes(),
	})
}

func (m *HealerRepository) GetTokenAddressBySymbol(ctx context.Context, chain string, tokenSymbol string) (*common.Address, error) {
	address, err := m.Queries.GetTokenAddressBySymbol(ctx, sqlc.GetTokenAddressBySymbolParams{
		Chain:  chain,
		Symbol: tokenSymbol,
	})
	if err != nil {
		tokens, _ := m.Queries.ListTokens(ctx)
		fmt.Printf("tokens: %v\n", tokens)
		return nil, err
	}
	addr := common.BytesToAddress(address)
	return &addr, nil
}
