package postgres

import (
	"context"

	"github.com/scalarorg/scalar-healer/pkg/db"
)

func (m *PostgresRepository) SaveUtxoSnapshot(ctx context.Context, utxoSnapshot *db.UTXOSnapshot) error {
	// filter := bson.M{"custodian_group_uid": utxoSnapshot.CustodianGroupUID}

	// update := mongo.Pipeline{
	// 	{{"$set": bson.M{
	// 		"block_height": bson.M{
	// 			"$cond": bson.M{
	// 				"if":   bson.M{"$gt": bson.A{"$$newHeight", "$block_height"}},
	// 				"then": "$$newHeight",
	// 				"else": "$block_height",
	// 			},
	// 		},
	// 		"utxos": bson.M{
	// 			"$cond": bson.M{
	// 				"if":   bson.M{"$gt": bson.A{"$$newHeight", "$block_height"}},
	// 				"then": "$$newUtxos",
	// 				"else": "$utxos",
	// 			},
	// 		},
	// 	}}},
	// }

	// updateOptions := options.Update().SetUpsert(true)

	// _, err := m.UtxoSnapshots.UpdateOne(
	// 	ctx,
	// 	filter,
	// 	update,
	// 	updateOptions,
	// 	options.UpdateOptions{
	// 		Let: bson.M{
	// 			"newHeight": utxoSnapshot.BlockHeight,
	// 			"newUtxos":  utxoSnapshot.UTXOs,
	// 		},
	// 	},
	// )

	// _, err = m.UtxoSnapshots.UpdateOne(ctx, filter, update)
	// if err != nil {
	// 	log.Error().Err(err).Msg("Failed to update existing UTXO snapshot")
	// 	return err
	// }
	return nil
}
