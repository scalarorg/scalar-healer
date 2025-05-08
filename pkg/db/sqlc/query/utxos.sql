-- name: SaveUTXOs :exec
INSERT INTO utxos (tx_id, vout, script_pubkey, amount_in_sats, custodian_group_uid, block_height)
VALUES (unnest($1::bytea[]), unnest($2::bigint[]), unnest($3::bytea[]), unnest($4::numeric[]), unnest($5::bytea[]), unnest($6::bigint[]));

-- name: GetUTXOs :many
SELECT * FROM utxos WHERE tx_id = ANY($1::bytea[]) AND vout = ANY($2::bigint[]);

-- name: GetUTXOsByCustodianGroupUID :many
SELECT * FROM utxos WHERE custodian_group_uid = $1::bytea;

-- name: DeleteUTXOs :exec
DELETE FROM utxos WHERE custodian_group_uid = $1::bytea;