package healer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) SaveTransferRequest(ctx context.Context, chain string, address common.Address, signature []byte, amount *big.Int, destChain string, destAddress *common.Address, symbol string, nonce uint64) error {
	return m.execTx(ctx, func(_ context.Context, q *sqlc.Queries) error {
		currentNonce := m.GetNonce(ctx, address)
		if nonce != currentNonce {
			return constants.ErrInvalidNonce
		}

		err := m.Queries.UpsertNonce(ctx, sqlc.UpsertNonceParams{
			Address: address.Bytes(),
			Nonce:   sqlc.ConvertUint64ToNumeric(currentNonce),
		})
		if err != nil {
			return err
		}

		return m.Queries.SaveTransferRequest(ctx, sqlc.SaveTransferRequestParams{
			Chain:              chain,
			Address:            address.Bytes(),
			Signature:          signature,
			Amount:             amount.String(),
			DestinationChain:   destChain,
			DestinationAddress: destAddress.Bytes(),
			Symbol:             symbol,
			Nonce:              sqlc.ConvertUint64ToNumeric(nonce),
		})
	})
}

func (m *HealerRepository) ListTransferRequests(ctx context.Context, address common.Address, page, size int32) ([]sqlc.TransferRequest, error) {
	return m.Queries.ListTransferRequests(ctx, sqlc.ListTransferRequestsParams{
		Address: address.Bytes(),
		Offset:  page * size,
		Limit:   size,
	})
}
