-- name: SaveChainRedeemSessions :exec
INSERT INTO chain_redeem_sessions (chain, custodian_group_uid, sequence, current_phase)
VALUES (unnest($1::text[]), unnest($2::bytea[]), unnest($3::bigint[]), unnest($4::text[])::redeem_phase)
ON CONFLICT (chain, custodian_group_uid) DO UPDATE
SET sequence = EXCLUDED.sequence, current_phase = EXCLUDED.current_phase;
