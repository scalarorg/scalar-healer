package postgres

import (
	"context"

	"github.com/scalarorg/scalar-healer/pkg/db"
)

func (m *PostgresRepository) GetAllCustodianGroups(ctx context.Context) ([]db.CustodianGroup, error) {
	// cursor, err := m.CustodianGroups.Find(ctx, bson.M{})
	// if err != nil {
	// 	return nil, err
	// }
	// defer cursor.Close(ctx)
	// var result []*CustodianGroup

	// if err := cursor.All(ctx, &result); err != nil {
	// 	return nil, err
	// }
	// return result, nil
	return nil, nil
}
