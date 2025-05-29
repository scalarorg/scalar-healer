-- name: SaveTransferRequest :exec
INSERT INTO transfer_requests (address,  signature, chain, destination_chain, destination_address, symbol, amount, nonce)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: ListTransferRequests :many
SELECT *
FROM transfer_requests
WHERE address = $1
ORDER BY nonce DESC
LIMIT $2 OFFSET $3;


