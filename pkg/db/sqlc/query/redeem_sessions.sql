-- name: SaveRedeemSessions :exec
INSERT INTO redeem_sessions (custodian_group_uid, sequence, current_phase, is_switching, phase_expired_at)
VALUES (unnest($1::bytea[]), unnest($2::bigint[]), unnest($3::text[])::redeem_phase, unnest($4::boolean[]), unnest($5::timestamp[]))
ON CONFLICT (custodian_group_uid) DO UPDATE
SET sequence = EXCLUDED.sequence, current_phase = EXCLUDED.current_phase, is_switching = EXCLUDED.is_switching, phase_expired_at = EXCLUDED.phase_expired_at;

-- name: GetRedeemSession :one
SELECT * FROM redeem_sessions WHERE custodian_group_uid = $1;

-- name: GetCompletedRedeemSessions :many
SELECT * FROM redeem_sessions WHERE is_switching = false;

-- name: SaveChainRedeemSessions :exec
INSERT INTO chain_redeem_sessions (chain, custodian_group_uid, sequence, current_phase)
VALUES (unnest($1::text[]), unnest($2::bytea[]), unnest($3::bigint[]), unnest($4::text[])::redeem_phase)
ON CONFLICT (chain, custodian_group_uid) DO UPDATE
SET sequence = EXCLUDED.sequence, current_phase = EXCLUDED.current_phase;

-- name: GetChainRedeemSession :one
SELECT * FROM chain_redeem_sessions WHERE custodian_group_uid = $1 AND chain = $2;