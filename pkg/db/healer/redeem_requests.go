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
	"github.com/scalarorg/bitcoin-vault/go-utils/encode"
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/evm"
	chains_utils "github.com/scalarorg/scalar-healer/pkg/utils/chains"
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

	err := m.execTx(ctx, func(txCtx context.Context, q *sqlc.Queries) error {
		protocol, err := m.GetProtocol(txCtx, symbol)
		if err != nil {
			return constants.ErrTokenNotExists
		}

		_, err = protocol.GetTokenDetailsByChain(sourceChain)
		if err != nil {
			return constants.ErrTokenNotExists
		}

		// redeem session
		redeemSession, err := m.GetRedeemSession(txCtx, protocol.CustodianGroupUid)
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

		requestId, _ := createRedeemCommandID(nonce, address, sourceChain, destChain, symbol, amount)

		reservedUtxos, newUtxos, err := m.reserveUtxos(txCtx, redeemSession.CustodianGroupUid, uint64(protocol.CustodianQuorum), requestId[:], amount.Uint64(), constants.CHAIN_PARAMS.RedeemTxsVsizeLimit)
		if err != nil {
			return err
		}


		err = m.saveUtxoSnapshot(txCtx, newUtxos)
		if err != nil {
			return err
		}

		redeemTokenParams := &evm.RedeemTokenPayloadWithType{
			RedeemTokenPayload: evm.RedeemTokenPayload{
				Amount:        amount.Uint64(),
				LockingScript: lockingScript,
				Utxos:         reservedUtxos,
				RequestId:     requestId,
			},
			PayloadType: encode.ContractCallWithTokenPayloadType_CustodianOnly,
		}

		params, err := redeemTokenParams.AbiPack()
		if err != nil {
			return err
		}

		chainID, err := chains_utils.ChainName(sourceChain).GetChainID()
		if err != nil {
			return err
		}

		redeemCommand, err := NewRedeemCommand(sourceChain, requestId, chainID, params)
		if err != nil {
			return err
		}

		err = m.Queries.SaveRedeemCommand(ctx, sqlc.SaveRedeemCommandParams{
			ID:      redeemCommand.ID,
			Chain:   redeemCommand.Chain,
			Status:  redeemCommand.Status,
			Params:  redeemCommand.Params,
			Data:    redeemCommand.Data,
			SigHash: redeemCommand.SigHash,
		})
		if err != nil {
			return err
		}

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

func createRedeemCommandID(nonce uint64, address common.Address, sourceChain, destChain, symbol string, amount *big.Int) ([32]byte, string) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, nonce)

	hash := crypto.Keccak256(bz, address.Bytes(), []byte(sourceChain), []byte(destChain), []byte(symbol), amount.Bytes())

	hashHex := hex.EncodeToString(hash)

	byte32Hash := [32]byte{}
	copy(byte32Hash[:], hash)

	return byte32Hash, hashHex
}

func (m *HealerRepository) reserveUtxos(ctx context.Context, grUid []byte, quorum uint64, requestID []byte, amount uint64, vsizeLimit uint64) ([]sqlc.Utxo, sqlc.UtxoSnapshot, error) {
	utxoSnapshot, err := m.GetUtxoSnapshot(ctx, grUid)
	if err != nil {
		return nil, nil, err
	}

	if len(utxoSnapshot) == 0 {
		return nil, nil, fmt.Errorf("no utxos found")
	}

	// // Find optimal UTXO combination using knapsack algorithm
	reserveUtxos, err := sqlc.UtxoSnapshot(utxoSnapshot).ReserveUtxos(requestID, amount, quorum, vsizeLimit)
	if err != nil {
		return nil, nil, err
	}

	return reserveUtxos, utxoSnapshot, nil
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
