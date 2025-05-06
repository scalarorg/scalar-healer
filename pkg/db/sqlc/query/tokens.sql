-- name: SaveTokens :exec
INSERT INTO tokens (address, chain_id, protocol, symbol, decimal, name, avatar, active)
VALUES (unnest($1::bytea[]), unnest($2::numeric[]), unnest($3::text[]), unnest($4::text[]), unnest($5::numeric[]), unnest($6::text[]), unnest($7::text[]), unnest($8::boolean[]));

-- name: GetTokenSymbolByAddress :one
SELECT symbol FROM tokens WHERE chain_id = $1 AND address = $2;

-- name: GetTokenAddressBySymbol :one
SELECT address FROM tokens WHERE chain_id = $1 AND symbol = $2;

-- name: ListTokens :many
SELECT * FROM tokens;