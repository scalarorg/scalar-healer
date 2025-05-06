package postgres

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *PostgresRepository) FindPendingRedeemsTransaction(ctx context.Context, chainId string, expectedConfirmBlock int32) ([]chains.RedeemTx, error) {
	return nil, nil
}
func (m *PostgresRepository) UpdateRedeemExecutedCommands(ctx context.Context, chainId string, txHashes []string) error {
	return nil
}

func (m *PostgresRepository) SaveRedeemTxs(ctx context.Context, redeemTxs []chains.RedeemTx) error {
	return nil
}

func (m *PostgresRepository) SaveRedeemRequest(ctx context.Context, chainId uint64, address common.Address, signature []byte, amount *big.Int, symbol string, nonce uint64) error {
	return m.Queries.SaveRedeemRequest(ctx, sqlc.SaveRedeemRequestParams{
		Address:   address.Bytes(),
		Signature: signature,
		Amount:    amount.String(),
		Symbol:    symbol,
		ChainID:   db.ConvertUint64ToNumeric(chainId),
		Nonce:     db.ConvertUint64ToNumeric(nonce),
	})
}

func (m *PostgresRepository) ListRedeemRequests(ctx context.Context, address common.Address, page, size int32) ([]sqlc.RedeemRequest, error) {
	return m.Queries.ListRedeemRequests(ctx, sqlc.ListRedeemRequestsParams{
		Address: address.Bytes(),
		Offset:  page * size,
		Limit:   size,
	})
}
