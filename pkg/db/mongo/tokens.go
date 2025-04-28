package mongo

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/pkg/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoRepository) SaveTokenInfos(ctx context.Context, tokens []models.Token) error {
	tokenDocs := make([]interface{}, len(tokens))
	tokenSymbols := bson.A{}
	for i, token := range tokens {
		tokenDocs[i] = token
		tokenSymbols = append(tokenSymbols, token.Symbol)
	}
	_, err := m.Tokens.DeleteMany(ctx, bson.M{
		"symbol": bson.M{
			"$in": tokenSymbols,
		},
	})
	_, err = m.Tokens.InsertMany(ctx, tokenDocs)
	return err
}

func (m *MongoRepository) GetTokenSymbolByAddress(ctx context.Context, chainId uint64, tokenAddress common.Address) (string, error) {

	filter := bson.M{
		"chain_id": chainId,
		"address":  tokenAddress.Bytes(),
	}
	var data struct {
		Symbol string `bson:"symbol"`
	}
	opts := options.FindOne().SetProjection(bson.M{
		"symbol": 1,
	})
	err := m.Tokens.FindOne(ctx, filter, opts).Decode(&data)
	if err != nil {
		return "", err
	}
	return data.Symbol, nil
}

func (m *MongoRepository) GetTokenAddressBySymbol(ctx context.Context, chainId uint64, tokenSymbol string) (*common.Address, error) {
	filter := bson.M{
		"chain_id": chainId,
		"symbol":   tokenSymbol,
	}
	var data struct {
		Address []byte `bson:"address"`
	}
	opts := options.FindOne().SetProjection(bson.M{
		"address": 1,
	})
	err := m.Tokens.FindOne(ctx, filter, opts).Decode(&data)
	if err != nil {
		return nil, err
	}
	add := common.BytesToAddress(data.Address)
	return &add, nil
}
