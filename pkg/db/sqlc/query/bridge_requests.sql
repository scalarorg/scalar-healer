-- name: SaveBridgeRequest :exec
INSERT INTO bridge_requests (address, signature, chain_id, tx_hash, nonce)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (chain_id, nonce) DO NOTHING;

-- name: ListBridgeRequests :many
SELECT *
FROM bridge_requests
WHERE address = $1
ORDER BY nonce DESC
LIMIT $2 OFFSET $3;