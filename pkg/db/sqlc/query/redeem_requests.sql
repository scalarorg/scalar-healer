-- name: SaveRedeemRequest :exec
INSERT INTO redeem_requests (address,  signature, chain, symbol, amount, nonce)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: ListRedeemRequests :many
SELECT *
FROM redeem_requests
WHERE address = $1
ORDER BY nonce DESC
LIMIT $2 OFFSET $3;


