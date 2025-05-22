package healer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) FindPendingRedeemsTransaction(ctx context.Context, chainId string, expectedConfirmBlock int32) ([]chains.RedeemTx, error) {
	return nil, nil
}
func (m *HealerRepository) UpdateRedeemExecutedCommands(ctx context.Context, chainId string, txHashes []string) error {
	return nil
}

func (m *HealerRepository) SaveRedeemTxs(ctx context.Context, redeemTxs []chains.RedeemTx) error {
	return nil
}

func (m *HealerRepository) SaveRedeemRequest(ctx context.Context, chainId uint64, address common.Address, signature []byte, amount *big.Int, symbol string, nonce uint64) error {
	return m.execTx(ctx, func(q *sqlc.Queries) error {
		currentNonce := m.GetNonce(ctx, address)
		if nonce != currentNonce {
			return constants.ErrInvalidNonce
		}

		err := m.Queries.UpsertNonce(ctx, sqlc.UpsertNonceParams{
			Address: address.Bytes(),
			Nonce:   db.ConvertUint64ToNumeric(currentNonce),
		})
		if err != nil {
			return err
		}

		return m.Queries.SaveRedeemRequest(ctx, sqlc.SaveRedeemRequestParams{
			Address:   address.Bytes(),
			Signature: signature,
			Amount:    amount.String(),
			Symbol:    symbol,
			ChainID:   db.ConvertUint64ToNumeric(chainId),
			Nonce:     db.ConvertUint64ToNumeric(nonce),
		})
	})
}

func (m *HealerRepository) ListRedeemRequests(ctx context.Context, address common.Address, page, size int32) ([]sqlc.RedeemRequest, error) {
	return m.Queries.ListRedeemRequests(ctx, sqlc.ListRedeemRequestsParams{
		Address: address.Bytes(),
		Offset:  page * size,
		Limit:   size,
	})
}
