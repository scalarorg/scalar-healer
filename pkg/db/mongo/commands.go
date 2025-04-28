package mongo

import (
	"context"

	"github.com/scalarorg/data-models/chains"
)

func (m *MongoRepository) UpdateEvmCommandExecuted(ctx context.Context, cmdExecuted *chains.CommandExecuted) error {
	return nil
}
