-- name: SaveTokens :exec
INSERT INTO tokens (address, chain, chain_id, symbol, active)
VALUES (unnest($1::bytea[]), unnest($2::text[]), unnest($3::numeric[]), unnest($4::text[]), unnest($5::boolean[]))
ON CONFLICT (symbol, chain) DO NOTHING;

-- name: GetTokenSymbolByAddress :one
SELECT symbol FROM tokens WHERE chain = $1 AND address = $2;

-- name: GetTokenAddressBySymbol :one
SELECT address FROM tokens WHERE chain = $1 AND symbol = $2;

-- name: ListTokens :many
SELECT * FROM tokens;