-- name: SaveRedeemRequest :one
INSERT INTO redeem_requests (address, source_chain, dest_chain, symbol, amount, locking_script, custodian_group_uid)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: ListRedeemRequests :many
SELECT rq.*, rc.status as status, rc.signature, rc.sig_hash, COUNT(rq.id) AS count
FROM redeem_requests rq
LEFT JOIN redeem_commands rc ON rq.id = rc.request_id
WHERE rq.address = $1
GROUP BY rq.id, rc.status, rc.signature, rc.sig_hash
ORDER BY rq.created_at DESC
LIMIT $2 OFFSET $3;
