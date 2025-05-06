
-- name: SaveTokens :exec
INSERT INTO tokens (address, chain_id, protocol, symbol, decimal, name, avatar)
VALUES (unnest($1::bytea[]), unnest($2::numeric[]), unnest($3::text[]), unnest($4::text[]), unnest($5::numeric[]), unnest($6::text[]), unnest($7::text[]));
