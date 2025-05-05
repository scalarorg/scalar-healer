package postgres

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *PostgresRepository) SaveBridgeRequest(ctx context.Context, chainId uint64, address common.Address, signature []byte, txHash []byte, nonce uint64) error {
	return m.Queries.SaveBridgeRequest(ctx, sqlc.SaveBridgeRequestParams{
		Address:   address.Bytes(),
		TxHash:    txHash,
		Signature: signature,
		ChainID:   db.ConvertUint64ToNumeric(chainId),
		Nonce:     db.ConvertUint64ToNumeric(nonce),
	})
}

func (m *PostgresRepository) ListBridgeRequests(ctx context.Context, address common.Address, page, size int32) ([]sqlc.BridgeRequest, error) {
	return m.Queries.ListBridgeRequests(ctx, sqlc.ListBridgeRequestsParams{
		Address: address.Bytes(),
		Limit:   size,
		Offset:  page * size,
	})
}
