package postgres

import (
	"context"

	"github.com/scalarorg/data-models/chains"
)

func (m *PostgresRepository) UpdateEvmCommandExecuted(ctx context.Context, cmdExecuted *chains.CommandExecuted) error {
	return nil
}
