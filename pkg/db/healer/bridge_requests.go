package healer

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) SaveBridgeRequest(ctx context.Context, chain string, address common.Address, signature []byte, txHash []byte, nonce uint64) error {
	// TODO: verify chain
	return m.execTx(ctx, func(cv context.Context, q *sqlc.Queries) error {
		currentNonce := m.GetNonce(ctx, address)
		if nonce != currentNonce {
			return constants.ErrInvalidNonce
		}

		err := m.Queries.UpsertNonce(cv, sqlc.UpsertNonceParams{
			Address: address.Bytes(),
			Nonce:   sqlc.ConvertUint64ToNumeric(currentNonce),
		})
		if err != nil {
			return err
		}

		return m.Queries.SaveBridgeRequest(cv, sqlc.SaveBridgeRequestParams{
			Address:   address.Bytes(),
			TxHash:    txHash,
			Signature: signature,
			Chain:     chain,
			Nonce:     sqlc.ConvertUint64ToNumeric(nonce),
		})
	})
}

func (m *HealerRepository) ListBridgeRequests(ctx context.Context, address common.Address, page, size int32) ([]sqlc.BridgeRequest, error) {
	return m.Queries.ListBridgeRequests(ctx, sqlc.ListBridgeRequestsParams{
		Address: address.Bytes(),
		Limit:   size,
		Offset:  page * size,
	})
}
