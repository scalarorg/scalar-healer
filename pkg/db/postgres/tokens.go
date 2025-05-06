package postgres

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *PostgresRepository) SaveTokens(ctx context.Context, tokens []sqlc.Token) error {

	var addresses [][]byte
	var chainIds []pgtype.Numeric
	var protocols []string
	var symbols []string
	var decimals []pgtype.Numeric
	var names []string
	var avatars []string

	for _, token := range tokens {
		addresses = append(addresses, token.Address)
		chainIds = append(chainIds, token.ChainID)
		protocols = append(protocols, token.Protocol)
		symbols = append(symbols, token.Symbol)
		decimals = append(decimals, token.Decimal)
		names = append(names, token.Name)
		avatars = append(avatars, token.Avatar)
	}

	return m.Queries.SaveTokens(ctx, sqlc.SaveTokensParams{
		Column1: addresses, // address
		Column2: chainIds,  // chain_id,
		Column3: protocols, // protocol
		Column4: symbols,   // symbol
		Column5: decimals,  // decimal
		Column6: names,     // name
		Column7: avatars,   // avatar
	})
}

func (m *PostgresRepository) GetTokenSymbolByAddress(ctx context.Context, chainId uint64, tokenAddress common.Address) (string, error) {

	// filter := bson.M{
	// 	"chain_id": chainId,
	// 	"address":  tokenAddress.Bytes(),
	// }
	// var data struct {
	// 	Symbol string `bson:"symbol"`
	// }
	// opts := options.FindOne().SetProjection(bson.M{
	// 	"symbol": 1,
	// })
	// err := m.Tokens.FindOne(ctx, filter, opts).Decode(&data)
	// if err != nil {
	// 	return "", err
	// }
	// return data.Symbol, nil
	return "", nil
}

func (m *PostgresRepository) GetTokenAddressBySymbol(ctx context.Context, chainId uint64, tokenSymbol string) (*common.Address, error) {
	// filter := bson.M{
	// 	"chain_id": chainId,
	// 	"symbol":   tokenSymbol,
	// }
	// var data struct {
	// 	Address []byte `bson:"address"`
	// }
	// opts := options.FindOne().SetProjection(bson.M{
	// 	"address": 1,
	// })
	// err := m.Tokens.FindOne(ctx, filter, opts).Decode(&data)
	// if err != nil {
	// 	return nil, err
	// }
	// add := common.BytesToAddress(data.Address)
	// return &add, nil
	return nil, nil
}
