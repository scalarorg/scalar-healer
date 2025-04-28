package mongo

import (
	"context"
	"encoding/hex"

	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongoRepository) GetAllCustodianGroups(ctx context.Context) ([]string, error) {
	protocolCollection := m.DB.Collection(COLLECTION_PROTOCOLS)
	cursor, err := protocolCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var result []string
	for cursor.Next(ctx) {
		var d struct {
			CustodianGroupUid [32]byte `bson:"custodian_group_uid"`
		}
		err := cursor.Decode(&d)
		if err != nil {
			return nil, err
		}
		result = append(result, hex.EncodeToString(d.CustodianGroupUid[:]))
	}
	return result, nil
}
