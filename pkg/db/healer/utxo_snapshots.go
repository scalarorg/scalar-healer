package healer

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jinzhu/copier"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) GetUtxoSnapshot(ctx context.Context, uid []byte) (sqlc.UtxoSnapshot, error) {

	utxoSnapshot, err := m.Queries.GetUtxoSnapshot(ctx, uid)
	if err != nil {
		return nil, err
	}

	utxos := make([]*sqlc.UtxoWithReservations, len(utxoSnapshot))
	for i, utxo := range utxoSnapshot {
		utxos[i] = &sqlc.UtxoWithReservations{}
		copier.Copy(utxos[i], &utxo)

		reservationsJson, err := json.Marshal(utxo.Reservations)
		if err != nil {
			log.Error().Err(err).Msgf("failed to marshal reservations for utxo %x", utxo.TxID)
			return nil, err
		}
		var reservationsRaw []struct {
			RequestID string `json:"request_id"`
			Amount    uint64 `json:"amount"`
		}
		err = json.Unmarshal(reservationsJson, &reservationsRaw)
		if err != nil {
			log.Error().Err(err).Msgf("failed to unmarshal reservations for utxo %x", utxo.TxID)
			return nil, err
		}

		var reservations []*sqlc.UtxoReservation
		for _, r := range reservationsRaw {
			var reqID []byte
			if len(r.RequestID) >= 2 && r.RequestID[:2] == "\\x" {
				reqID, err = hex.DecodeString(r.RequestID[2:])
				if err != nil {
					log.Error().Err(err).Msgf("failed to decode request_id %s", r.RequestID)
					return nil, err
				}
			}
			reservations = append(reservations, &sqlc.UtxoReservation{
				ReservationID: reqID,
				Amount:        sqlc.ConvertUint64ToNumeric(r.Amount),
				UtxoTxID:      utxo.TxID,
				UtxoVout:      utxo.Vout,
			})
		}
		utxos[i].Reservations = reservations

	}
	return utxos, nil
}

type UtxoReservationWithRequest struct {
	UtxoTxID  []byte
	UtxoVout  int64
	RequestID string
}

func (m *HealerRepository) SaveUtxoSnapshot(ctx context.Context, utxoSnapshot sqlc.UtxoSnapshot) error {
	return m.execTx(ctx, func(cv context.Context, q *sqlc.Queries) error { return m.saveUtxoSnapshot(cv, utxoSnapshot) })
}

func (m *HealerRepository) saveUtxoSnapshot(ctx context.Context, utxoSnapshot sqlc.UtxoSnapshot) error {
	return m.requireTx(ctx, func(cv context.Context) error {
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

			reservationsMap map[string]sqlc.Reservation = make(map[string]sqlc.Reservation)

			utxoReservations []sqlc.UtxoReservation = make([]sqlc.UtxoReservation, 0)
		)

		var firstGrUID = utxoSnapshot[0].CustodianGroupUid
		var firstBlockHeight = utxoSnapshot[0].BlockHeight
		var firstScriptPubkey = utxoSnapshot[0].ScriptPubkey

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

			if len(utxo.Reservations) == 0 {
				continue
			}

			for _, reservation := range utxo.Reservations {
				if _, ok := reservationsMap[hex.EncodeToString(reservation.ReservationID)]; !ok {
					reservationsMap[hex.EncodeToString(reservation.ReservationID)] = sqlc.Reservation{
						RequestID: reservation.ReservationID,
					}
				}

				utxoReservations = append(utxoReservations, sqlc.UtxoReservation{
					UtxoTxID:      utxo.TxID,
					UtxoVout:      utxo.Vout,
					ReservationID: reservation.ReservationID,
					Amount:        reservation.Amount,
				})
			}
		}

		utxos, err := m.Queries.GetUTXOsByCustodianGroupUID(cv, firstGrUID)
		if err != nil && err != pgx.ErrNoRows {
			return err
		}

		if len(utxos) > 0 {
			if firstBlockHeight < utxos[0].BlockHeight {
				return fmt.Errorf("UTXO snapshot is not newer than the existing one")
			}
		}

		// TODO: Currently, we delete all utxos then insert again, consider updating instead

		err = m.deleteUtxoSnapshot(cv, firstGrUID)
		if err != nil {
			log.Error().Err(err).Msgf("failed to delete existing UTXO snapshot for custodian group %s", firstGrUID)
			return err
		}

		// insert new UTXO snapshot
		err = m.Queries.SaveUTXOs(cv, sqlc.SaveUTXOsParams{
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

		requestIDs := make([][]byte, 0)
		for _, reservation := range reservationsMap {
			requestIDs = append(requestIDs, reservation.RequestID)
		}

		_, err = m.Queries.SaveReservations(cv, requestIDs)
		if err != nil {
			return err
		}

		var (
			reservationTxIDs [][]byte
			reservationVouts []int64
			reservationIDs   [][]byte
			amounts          []pgtype.Numeric
		)

		for _, utxoReservation := range utxoReservations {
			reservationTxIDs = append(reservationTxIDs, utxoReservation.UtxoTxID)
			reservationVouts = append(reservationVouts, utxoReservation.UtxoVout)
			reservationIDs = append(reservationIDs, utxoReservation.ReservationID)
			amounts = append(amounts, utxoReservation.Amount)
		}

		return m.Queries.SaveUtxoReservations(cv, sqlc.SaveUtxoReservationsParams{
			Column1: reservationTxIDs,
			Column2: reservationVouts,
			Column3: reservationIDs,
			Column4: amounts,
		})
	})
}

func (m *HealerRepository) deleteUtxoSnapshot(ctx context.Context, grUID []byte) error {
	return m.requireTx(ctx, func(cv context.Context) error {
		log.Info().Msgf("deleting utxo snapshot for custodian group %x", grUID)
		err := m.Queries.DeleteUTXOs(cv, grUID)
		if err != nil {
			return err
		}

		return m.Queries.DeleteReservations(cv)
	})
}
