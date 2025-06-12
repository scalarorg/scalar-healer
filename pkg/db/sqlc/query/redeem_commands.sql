-- name: SaveRedeemCommands :exec
INSERT INTO redeem_commands (id, chain, status, params, data, sig_hash, request_id)
VALUES (unnest($1::bytea[]), unnest($2::text[]), unnest($3::text[])::batch_status, unnest($4::bytea[]), unnest($5::bytea[]), unnest($6::bytea[]), unnest($7::bigint[]));

-- name: SaveRedeemCommand :exec
INSERT INTO redeem_commands (id, chain, status, params, data, sig_hash, request_id)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: ListPendingSigningRedeemCommands :many
SELECT * FROM redeem_commands WHERE status = 'PENDING' AND signature IS NULL;

-- name: SubmitRedeemCommandSignature :exec
UPDATE redeem_commands SET signature = $1, status = 'SIGNED' WHERE id = $2;
