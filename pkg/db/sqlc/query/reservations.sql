-- name: SaveReservations :many
INSERT INTO reservations (
    request_id,
    amount
) VALUES (
    unnest($1::text[]),
    unnest($2::numeric[])
) RETURNING id, request_id;

-- name: DeleteReservations :exec
DELETE FROM reservations WHERE id IN (
  SELECT r.id
  FROM reservations r
  LEFT JOIN utxo_reservations ur ON r.id = ur.reservation_id
  WHERE ur.reservation_id IS NULL
);