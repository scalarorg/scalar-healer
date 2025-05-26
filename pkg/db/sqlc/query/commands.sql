-- name: SaveCommands :exec
INSERT INTO commands (command_id, chain, params, status, command_type)
VALUES (unnest($1::bytea[]), unnest($2::text[]), unnest($3::bytea[]), unnest($4::int[]), unnest($5::command_type[]));