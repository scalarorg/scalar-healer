package mongo

import (
	"context"

	"github.com/scalarorg/scalar-healer/pkg/db/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongoRepository) GetAllCustodianGroups(ctx context.Context) ([]*models.CustodianGroup, error) {
	cursor, err := m.CustodianGroups.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var result []*models.CustodianGroup

	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}
