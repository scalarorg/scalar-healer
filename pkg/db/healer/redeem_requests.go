package healer

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/copier"
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/scalar-healer/constants"
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

	err := m.execTx(ctx, func(q *sqlc.Queries) error {
		protocol, err := m.GetProtocol(ctx, symbol)
		if err != nil {
			return constants.ErrTokenNotExists
		}

		_, err = protocol.GetTokenDetailsByChain(sourceChain)
		if err != nil {
			return constants.ErrTokenNotExists
		}

		// redeem session
		redeemSession, err := m.GetRedeemSession(ctx, protocol.CustodianGroupUid)
		if err != nil {
			return constants.ErrInvalidRedeemSession
		}

		if redeemSession.IsSwitching.Bool {
			return constants.ErrRedeemSessionSwitching
		}

		if redeemSession.CurrentPhase == sqlc.RedeemPhaseEXECUTING {
			// TODO: handle for executing phase
			return errors.New("executing phase is not supported yet")
		}

		// phase is preparing

		// 1. Create a reservation
		// 2. Store the reservation
		// 3. Create cmd to sign on it
		// 4. Save the cmd

		return m.Queries.SaveRedeemRequest(ctx, sqlc.SaveRedeemRequestParams{
			Address:           address.Bytes(),
			Amount:            amount.String(),
			Symbol:            symbol,
			SourceChain:       sourceChain,
			DestChain:         destChain,
			LockingScript:     lockingScript,
			CustodianGroupUid: redeemSession.CustodianGroupUid,
		})
	})

	return err

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
