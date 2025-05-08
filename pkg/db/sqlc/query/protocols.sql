-- name: GetProtocol :one
SELECT * FROM protocols WHERE asset = $1;

-- name: SaveProtocols :exec
INSERT INTO protocols (asset, name, custodian_group_name, custodian_group_uid, tag, liquidity_model, symbol, decimals, capacity, daily_mint_limit, avatar)
VALUES(unnest($1::text[]), unnest($2::text[]), unnest($3::text[]), unnest($4::bytea[]), unnest($5::text[]), unnest($6::text[]), unnest($7::text[]), unnest($8::bigint[]), unnest($9::numeric[]), unnest($10::numeric[]), unnest($11::text[]));
