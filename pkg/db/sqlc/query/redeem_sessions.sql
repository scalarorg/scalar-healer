-- name: SaveRedeemSessions :exec
INSERT INTO redeem_sessions (custodian_group_uid, sequence, current_phase)
VALUES (unnest($1::bytea[]), unnest($2::bigint[]), unnest($3::text[])::redeem_phase);
