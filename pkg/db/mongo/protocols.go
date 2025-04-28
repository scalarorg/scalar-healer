package mongo

import (
	"context"

	"github.com/scalarorg/scalar-healer/pkg/db/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongoRepository) SaveProtocols(ctx context.Context, protocols []models.Protocol) error {
	protocolDocs := make([]interface{}, len(protocols))
	protocolNames := bson.A{}
	for i, protocol := range protocols {
		protocolDocs[i] = protocol
		protocolNames = append(protocolNames, protocol.Name)
	}
	filter := bson.M{
		"name": bson.M{
			"$in": protocolNames,
		},
	}
	collection := m.DB.Collection(COLLECTION_PROTOCOLS)
	collection.DeleteMany(ctx, filter)
	_, err := collection.InsertMany(ctx, protocolDocs)
	return err
}
