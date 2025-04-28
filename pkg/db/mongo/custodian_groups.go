package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoRepository) GetAllCustodianGroups(ctx context.Context) ([]string, error) {
	protocolCollection := m.DB.Collection(COLLECTION_PROTOCOLS)

	options := options.Find().SetProjection(bson.M{
		"custodian_group_uid": 1,
	})

	cursor, err := protocolCollection.Find(ctx, bson.M{}, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var result []string

	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}
