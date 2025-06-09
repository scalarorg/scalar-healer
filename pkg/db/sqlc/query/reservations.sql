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

-- name: DeleteReservations :exec
DELETE FROM reservations WHERE id IN (
  SELECT r.id
  FROM reservations r
  LEFT JOIN utxo_reservations ur ON r.id = ur.reservation_id
  WHERE ur.reservation_id IS NULL
);