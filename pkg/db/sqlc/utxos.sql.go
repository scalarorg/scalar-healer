// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: utxos.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deleteUTXOs = `-- name: DeleteUTXOs :exec
DELETE FROM utxos WHERE custodian_group_uid = $1::bytea
`

func (q *Queries) DeleteUTXOs(ctx context.Context, dollar_1 []byte) error {
	_, err := q.db.Exec(ctx, deleteUTXOs, dollar_1)
	return err
}

const getUTXOs = `-- name: GetUTXOs :many
SELECT id, tx_id, vout, script_pubkey, amount_in_sats, custodian_group_uid, block_height, created_at, updated_at FROM utxos WHERE tx_id = ANY($1::bytea[]) AND vout = ANY($2::bigint[])
`

type GetUTXOsParams struct {
	Column1 [][]byte `json:"column_1"`
	Column2 []int64  `json:"column_2"`
}

func (q *Queries) GetUTXOs(ctx context.Context, arg GetUTXOsParams) ([]Utxo, error) {
	rows, err := q.db.Query(ctx, getUTXOs, arg.Column1, arg.Column2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Utxo{}
	for rows.Next() {
		var i Utxo
		if err := rows.Scan(
			&i.ID,
			&i.TxID,
			&i.Vout,
			&i.ScriptPubkey,
			&i.AmountInSats,
			&i.CustodianGroupUid,
			&i.BlockHeight,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUTXOsByCustodianGroupUID = `-- name: GetUTXOsByCustodianGroupUID :many
SELECT id, tx_id, vout, script_pubkey, amount_in_sats, custodian_group_uid, block_height, created_at, updated_at FROM utxos WHERE custodian_group_uid = $1::bytea
`

func (q *Queries) GetUTXOsByCustodianGroupUID(ctx context.Context, dollar_1 []byte) ([]Utxo, error) {
	rows, err := q.db.Query(ctx, getUTXOsByCustodianGroupUID, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Utxo{}
	for rows.Next() {
		var i Utxo
		if err := rows.Scan(
			&i.ID,
			&i.TxID,
			&i.Vout,
			&i.ScriptPubkey,
			&i.AmountInSats,
			&i.CustodianGroupUid,
			&i.BlockHeight,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveUTXOs = `-- name: SaveUTXOs :exec
INSERT INTO utxos (tx_id, vout, script_pubkey, amount_in_sats, custodian_group_uid, block_height)
VALUES (unnest($1::bytea[]), unnest($2::bigint[]), unnest($3::bytea[]), unnest($4::numeric[]), unnest($5::bytea[]), unnest($6::bigint[]))
`

type SaveUTXOsParams struct {
	Column1 [][]byte         `json:"column_1"`
	Column2 []int64          `json:"column_2"`
	Column3 [][]byte         `json:"column_3"`
	Column4 []pgtype.Numeric `json:"column_4"`
	Column5 [][]byte         `json:"column_5"`
	Column6 []int64          `json:"column_6"`
}

func (q *Queries) SaveUTXOs(ctx context.Context, arg SaveUTXOsParams) error {
	_, err := q.db.Exec(ctx, saveUTXOs,
		arg.Column1,
		arg.Column2,
		arg.Column3,
		arg.Column4,
		arg.Column5,
		arg.Column6,
	)
	return err
}
