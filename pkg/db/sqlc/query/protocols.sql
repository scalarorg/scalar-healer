-- name: GetProtocol :many
SELECT
    p.*,
    cg.custodians,
    cg.quorum,
    t.chain,
    t.chain_id,
    t.address
FROM protocols p
JOIN custodian_groups cg ON cg.uid = p.custodian_group_uid
LEFT JOIN tokens t ON t.symbol = p.symbol
WHERE p.symbol = $1
GROUP BY 
    p.id,
    p.symbol,
    p.name,
    p.bitcoin_pubkey,
    p.custodian_group_name,
    p.custodian_group_uid,
    cg.custodians,
    cg.quorum,
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

-- name: GetProtocols :many
SELECT
    p.*,
    cg.custodians,
    cg.quorum,
    t.chain,
    t.chain_id,
    t.address
FROM protocols p
JOIN custodian_groups cg ON cg.uid = p.custodian_group_uid
LEFT JOIN tokens t ON t.symbol = p.symbol
GROUP BY 
    p.id,
    p.symbol,
    p.name,
    p.bitcoin_pubkey,
    p.custodian_group_name,
    p.custodian_group_uid,
    cg.custodians,
    cg.quorum,
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
INSERT INTO protocols (symbol, name, bitcoin_pubkey, custodian_group_name, custodian_group_uid, tag, liquidity_model, decimals, capacity, daily_mint_limit, avatar)
VALUES(unnest($1::text[]), unnest($2::text[]), unnest($3::bytea[]), unnest($4::text[]), unnest($5::bytea[]), unnest($6::text[]), unnest($7::text[]), unnest($8::bigint[]), unnest($9::numeric[]), unnest($10::numeric[]), unnest($11::text[])) ON CONFLICT DO NOTHING;