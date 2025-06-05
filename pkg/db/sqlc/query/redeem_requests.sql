-- name: SaveRedeemRequest :exec
INSERT INTO redeem_requests (address,  source_chain, dest_chain, symbol, amount, locking_script, custodian_group_uid)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: ListRedeemRequests :many
SELECT *, COUNT(*) OVER() AS count
FROM redeem_requests
WHERE address = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
