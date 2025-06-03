package healer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/copier"
	"github.com/scalarorg/data-models/chains"
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

func (m *HealerRepository) SaveRedeemRequest(ctx context.Context, sourceChain, destChain string, address common.Address, amount *big.Int, symbol string, lockingScript []byte) error {
	return m.Queries.SaveRedeemRequest(ctx, sqlc.SaveRedeemRequestParams{
		Address:       address.Bytes(),
		Amount:        amount.String(),
		Symbol:        symbol,
		SourceChain:   sourceChain,
		DestChain:     destChain,
		LockingScript: lockingScript,
	})
}

func (m *HealerRepository) ListRedeemRequests(ctx context.Context, address common.Address, page, size int32) ([]sqlc.RedeemRequest, int64, error) {
	result, err := m.Queries.ListRedeemRequests(ctx, sqlc.ListRedeemRequestsParams{
		Address: address.Bytes(),
		Offset:  page * size,
		Limit:   size,
	})
	if err != nil {
		return nil, 0, err
	}

	if len(result) == 0 {
		return nil, 0, nil
	}

	var redeemRequests []sqlc.RedeemRequest
	for _, redeemRequest := range result {
		var req sqlc.RedeemRequest
		copier.Copy(&req, &redeemRequest)
		redeemRequests = append(redeemRequests, req)
	}
	return redeemRequests, result[0].Count, nil
}
