package mongo

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/pkg/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoRepository) SaveTokenInfos(ctx context.Context, tokens []models.Token) error {
	if len(tokens) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, len(tokens))
	for i, token := range tokens {
		filter := bson.M{"symbol": token.Symbol}
		update := bson.M{"$set": token}
		models[i] = mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
	}

	_, err := m.Tokens.BulkWrite(ctx, models)
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
