-- name: SaveRedeemCommands :exec
INSERT INTO redeem_commands (id, chain, status, params, data, sig_hash)
VALUES (unnest($1::bytea[]), unnest($2::text[]), unnest($3::int[]), unnest($4::bytea[]), unnest($5::bytea[]), unnest($6::bytea[]));

-- name: SaveRedeemCommand :exec
INSERT INTO redeem_commands (id, chain, status, params, data, sig_hash)
VALUES ($1, $2, $3, $4, $5, $6);