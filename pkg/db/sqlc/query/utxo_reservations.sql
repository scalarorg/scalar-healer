-- name: SaveUtxoReservations :exec
INSERT INTO utxo_reservations (utxo_tx_id, utxo_vout, reservation_id)
VALUES (unnest($1::bytea[]), unnest($2::bigint[]), unnest($3::bigint[])); 
