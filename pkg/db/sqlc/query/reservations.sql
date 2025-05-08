-- name: SaveReservations :exec
INSERT INTO reservations (
    utxo_tx_id,
    utxo_vout,
    request_id,
    amount
) VALUES (
    unnest($1::bytea[]),
    unnest($2::bigint[]),
    unnest($3::text[]),
    unnest($4::numeric[])
);