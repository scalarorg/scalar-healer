-- name: SaveBridgeRequest :exec
INSERT INTO bridge_requests (address, signature, chain, tx_hash, nonce)
VALUES ($1, $2, $3, $4, $5);

-- name: ListBridgeRequests :many
SELECT *
FROM bridge_requests
WHERE address = $1
ORDER BY nonce DESC
LIMIT $2 OFFSET $3;

