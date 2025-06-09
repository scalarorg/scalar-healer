package healer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jinzhu/copier"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) GetUtxoSnapshot(ctx context.Context, uid []byte) ([]sqlc.UtxoWithReservations, error) {
	utxoSnapshot, err := m.Queries.GetUtxoSnapshot(ctx, uid)
	if err != nil {
		return nil, err
	}

	utxos := make([]sqlc.UtxoWithReservations, len(utxoSnapshot))
	for i, utxo := range utxoSnapshot {
		utxos[i] = sqlc.UtxoWithReservations{}
		copier.Copy(&utxos[i], &utxo)

		var reservations []sqlc.Reservation

		err := json.Unmarshal(utxo.Reservations, &reservations)
		if err != nil {
			return nil, err
		}
		utxos[i].Reservations = reservations
	}
	return utxos, nil
}

func (m *HealerRepository) SaveUtxoSnapshot(ctx context.Context, utxoSnapshot []sqlc.UtxoWithReservations) error {

	return m.execTx(ctx, func(cv context.Context, q *sqlc.Queries) error {
		if len(utxoSnapshot) == 0 {
			return nil
		}

		var (
			txIDs              [][]byte
			vouts              []int64
			scriptPubkeys      [][]byte
			amountsInSats      []pgtype.Numeric
			custodianGroupUIDs [][]byte
			blockHeights       []int64

			requestIDs []string
			amounts    []pgtype.Numeric

			utxoReservationsMap map[string]struct {
				txID []byte
				vout int64
			}
		)

		var firstGrUID = utxoSnapshot[0].CustodianGroupUid
		var firstBlockHeight = utxoSnapshot[0].BlockHeight
		var firstScriptPubkey = utxoSnapshot[0].ScriptPubkey

		utxoReservationsMap = make(map[string]struct {
			txID []byte
			vout int64
		})

		for _, utxo := range utxoSnapshot {
			if !bytes.Equal(firstScriptPubkey, utxo.ScriptPubkey) {
				return fmt.Errorf("UTXO snapshot is not owned by the same script pubkey")
			}
			if firstBlockHeight != utxo.BlockHeight {
				return fmt.Errorf("UTXO snapshot is not owned by the same block height")
			}
			if !bytes.Equal(firstGrUID, utxo.CustodianGroupUid) {
				return fmt.Errorf("UTXO snapshot is not owned by the same custodian group")
			}

			txIDs = append(txIDs, utxo.TxID)
			vouts = append(vouts, utxo.Vout)
			scriptPubkeys = append(scriptPubkeys, utxo.ScriptPubkey)
			amountsInSats = append(amountsInSats, utxo.AmountInSats)
			custodianGroupUIDs = append(custodianGroupUIDs, utxo.CustodianGroupUid)
			blockHeights = append(blockHeights, utxo.BlockHeight)

			for _, reservation := range utxo.Reservations {
				requestIDs = append(requestIDs, reservation.RequestID)
				amounts = append(amounts, reservation.Amount)
				utxoReservationsMap[reservation.RequestID] = struct {
					txID []byte
					vout int64
				}{
					txID: utxo.TxID,
					vout: utxo.Vout,
				}
			}
		}

		utxos, err := q.GetUTXOsByCustodianGroupUID(ctx, firstGrUID)
		if err != nil && err != pgx.ErrNoRows {
			return err
		}

		if len(utxos) > 0 {
			if utxos[0].BlockHeight >= firstBlockHeight {
				return fmt.Errorf("UTXO snapshot is not newer than the existing one")
			}
		}

		err = m.deleteUtxoSnapshot(cv, firstGrUID)
		if err != nil {
			log.Error().Err(err).Msgf("failed to delete existing UTXO snapshot for custodian group %s", firstGrUID)
			return err
		}

		// insert new UTXO snapshot
		err = m.Queries.SaveUTXOs(ctx, sqlc.SaveUTXOsParams{
			Column1: txIDs,
			Column2: vouts,
			Column3: scriptPubkeys,
			Column4: amountsInSats,
			Column5: custodianGroupUIDs,
			Column6: blockHeights,
		})
		if err != nil {
			return err
		}

		reservationIDRows, err := m.Queries.SaveReservations(ctx, sqlc.SaveReservationsParams{
			Column1: requestIDs,
			Column2: amounts,
		})
		if err != nil {
			return err
		}

		var (
			reservationTxIDs [][]byte
			reservationVouts []int64
			reservationIDs   []int64
		)

		for _, row := range reservationIDRows {
			reservationTxIDs = append(reservationTxIDs, utxoReservationsMap[row.RequestID].txID)
			reservationVouts = append(reservationVouts, utxoReservationsMap[row.RequestID].vout)
			reservationIDs = append(reservationIDs, row.ID)
		}

		return m.Queries.SaveUtxoReservations(ctx, sqlc.SaveUtxoReservationsParams{
			Column1: reservationTxIDs,
			Column2: reservationVouts,
			Column3: reservationIDs,
		})
	})
}

func (m *HealerRepository) deleteUtxoSnapshot(ctx context.Context, grUID []byte) error {
	return m.requireTx(ctx, func(cv context.Context) error {

		log.Info().Msgf("deleting utxo snapshot for custodian group %s", grUID)

		err := m.Queries.DeleteUTXOs(cv, grUID)
		if err != nil {
			return err
		}

		return m.Queries.DeleteReservations(cv)
	})
}
