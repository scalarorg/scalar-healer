-- name: SaveReservations :many
INSERT INTO reservations (
    request_id
) VALUES (
    unnest($1::bytea[])
) RETURNING request_id;

-- name: DeleteReservations :exec
DELETE FROM reservations WHERE request_id IN (
  SELECT r.request_id
  FROM reservations r
  LEFT JOIN utxo_reservations ur ON r.request_id = ur.reservation_id
  WHERE ur.reservation_id IS NULL
);