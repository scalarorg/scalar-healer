package healer

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

		nonce := uint64(time.Now().UnixNano())

		hash := createRedeemCommandID(nonce, address, sourceChain, destChain, symbol, amount)

		hashHex := hex.EncodeToString(hash)

		byte32Hash := [32]byte{}
		copy(byte32Hash[:], hash)

		reservedUtxos, err := m.reserveUtxos(ctx, redeemSession.CustodianGroupUid, hashHex, amount.Uint64(), constants.CHAIN_PARAMS.RedeemTxsVsizeLimit)
		if err != nil {
			return err
		}

		_ = reservedUtxos

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

func createRedeemCommandID(nonce uint64, address common.Address, sourceChain, destChain, symbol string, amount *big.Int) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, nonce)

	return crypto.Keccak256(bz, address.Bytes(), []byte(sourceChain), []byte(destChain), []byte(symbol), amount.Bytes())
}

func (m *HealerRepository) reserveUtxos(ctx context.Context, grUid []byte, requestID string, amount uint64, vsizeLimit uint64) ([]sqlc.Utxo, error) {

	utxoSnapshot, err := m.GetUtxoSnapshot(ctx, grUid)
	if err != nil {
		return nil, err
	}

	if len(utxoSnapshot) == 0 {
		return nil, fmt.Errorf("no utxos found")
	}

	// // Find optimal UTXO combination using knapsack algorithm
	// reserveUtxos, err := utxoSnapshot.ReserveUtxos(requestID, amount, uint64(gr.Quorum), vsizeLimit)
	// if err != nil {
	// 	k.Logger(ctx).Error("failed to reserve utxos", "error", err)
	// 	return nil, err
	// }

	// // Update the redeem session in storage
	// k.setUtxoSnapshot(ctx, utxoSnapshot)

	// return reserveUtxos, nil

	return nil, nil
}

// func CreateRedeemParams(ctx context.Context, vsizeLimit uint64, amount *big.Int, sequence uint64, utxos []sqlc.Utxo) ([]byte, *CommandID, error) {

// 	dataHex := hex.EncodeToString(dataHash)

// 	// reservedUtxos, err := k.reserveUtxos(ctx, r.CustodianGroupUid, dataHex, amount,
// 	// 	vsizeLimit)
// 	// if err != nil {
// 	// 	return nil, nil, err
// 	// }

// 	cmdId := NewCommandID(dataHash, r.SourceChain)
// 	payload := &cov.RedeemTokenPayloadWithType{
// 		RedeemTokenPayload: cov.RedeemTokenPayload{
// 			Amount:        req.Amount,
// 			LockingScript: req.LockingScript,
// 			Utxos:         reservedUtxos,
// 			RequestId:     cmdId.Bytes(),
// 		},
// 		PayloadType: encode.ContractCallWithTokenPayloadType_CustodianOnly,
// 	}
// 	redeemTokenParams := &cov.RedeemTokenParams{
// 		DestinationChain:   req.DestChain.String(),
// 		DestinationAddress: req.Address,
// 		Payload:            *payload,
// 		Symbol:             req.Symbol,
// 		Amount:             req.Amount,
// 		CustodianGroupUID:  custodianGrUID,
// 		SessionSequence:    sequence,
// 	}
// 	params, err := redeemTokenParams.AbiPack()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return params, &cmdId, err
// }

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
