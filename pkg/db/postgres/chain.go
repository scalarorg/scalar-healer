package postgres

import "context"

func (m *PostgresRepository) GetChainName(ctx context.Context, chainType string, chainId uint64) (string, error) {
	return "", nil
}
