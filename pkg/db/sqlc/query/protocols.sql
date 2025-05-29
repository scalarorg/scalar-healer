-- name: GetProtocol :one
SELECT * FROM protocols WHERE symbol = $1;

-- name: GetProtocols :many
SELECT
    p.*, 
    t.chain,
    t.chain_id,
    t.address
FROM protocols p
LEFT JOIN tokens t ON t.symbol = p.symbol
GROUP BY 
    p.id,
    p.symbol,
    p.name,
    p.custodian_group_name,
    p.custodian_group_uid,
    p.tag,
    p.decimals, 
    p.liquidity_model,
    p.avatar,
    p.capacity,
    p.daily_mint_limit,
    p.created_at,
    p.updated_at,
    t.chain,
    t.chain_id,
    t.address;


-- name: SaveProtocols :exec
INSERT INTO protocols (symbol, name, custodian_group_name, custodian_group_uid, tag, liquidity_model, decimals, capacity, daily_mint_limit, avatar)
VALUES(unnest($1::text[]), unnest($2::text[]), unnest($3::text[]), unnest($4::bytea[]), unnest($5::text[]), unnest($6::text[]), unnest($7::bigint[]), unnest($8::numeric[]), unnest($9::numeric[]), unnest($10::text[])) ON CONFLICT DO NOTHING;