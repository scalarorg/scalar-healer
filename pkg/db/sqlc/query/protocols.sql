-- name: GetProtocol :one
SELECT * FROM protocols WHERE symbol = $1;

-- name: GetProtocols :many
SELECT * FROM protocols;

-- name: SaveProtocols :exec
INSERT INTO protocols (symbol, name, custodian_group_name, custodian_group_uid, tag, liquidity_model, decimals, capacity, daily_mint_limit, avatar)
VALUES(unnest($1::text[]), unnest($2::text[]), unnest($3::text[]), unnest($4::bytea[]), unnest($5::text[]), unnest($6::text[]), unnest($7::bigint[]), unnest($8::numeric[]), unnest($9::numeric[]), unnest($10::text[]));
