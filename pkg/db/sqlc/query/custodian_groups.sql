-- name: GetAllCustodianGroups :many
SELECT * FROM custodian_groups;

-- name: GetCustodianGroupByUID :one
SELECT * FROM custodian_groups WHERE uid = $1;

-- name: SaveCustodianGroups :exec
INSERT INTO custodian_groups (uid, name, bitcoin_pubkey, quorum)
VALUES (unnest($1::bytea[]), unnest($2::text[]), unnest($3::bytea[]), unnest($4::bigint[]));
