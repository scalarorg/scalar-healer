package healer

import (
	"bytes"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) SaveUtxoSnapshot(ctx context.Context, utxoSnapshot []sqlc.Utxo) error {
	return m.execTx(ctx, func(q *sqlc.Queries) error {
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
		}

		utxos, err := m.Queries.GetUTXOsByCustodianGroupUID(ctx, firstGrUID)
		if err != nil && err != pgx.ErrNoRows {
			return err
		}

		if len(utxos) > 0 {
			if utxos[0].BlockHeight >= firstBlockHeight {
				return fmt.Errorf("UTXO snapshot is not newer than the existing one")
			}
		}

		err = m.Queries.DeleteUTXOs(ctx, firstGrUID)
		if err != nil {
			return err
		}

		// insert new UTXO snapshot
		return m.Queries.SaveUTXOs(ctx, sqlc.SaveUTXOsParams{
			Column1: txIDs,
			Column2: vouts,
			Column3: scriptPubkeys,
			Column4: amountsInSats,
			Column5: custodianGroupUIDs,
			Column6: blockHeights,
		})
	})
}
