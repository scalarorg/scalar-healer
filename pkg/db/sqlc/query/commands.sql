-- name: SaveCommands :exec
INSERT INTO commands (id, chain, params, status, command_type, payload)
VALUES (unnest($1::bytea[]), unnest($2::text[]), unnest($3::bytea[]), unnest($4::int[]), unnest($5::text[])::command_type, unnest($6::bytea[]));

-- name: SaveCommandBatches :exec
INSERT INTO command_batchs (id, chain, data, sig_hash, status, extra_data)
VALUES (unnest($1::bytea[]), unnest($2::text[]), unnest($3::bytea[]), unnest($4::bytea[]), unnest($5::text[])::batch_status, unnest($6::bytea[]));

-- name: GetCommandBatches :many
SELECT * FROM command_batchs;

-- name: GetCommandBatchByID :one
SELECT * FROM command_batchs WHERE id = $1;