-- name: CreateGatewayAddress :exec
INSERT INTO gateway_addresses (address, chain_id)
VALUES ($1, $2);

-- name: CreateGatewayAddresses :exec
INSERT INTO gateway_addresses (address, chain_id)
VALUES (unnest($1::bytea[]), unnest($2::numeric[]));

-- name: GetGatewayAddress :one
SELECT address
FROM gateway_addresses
WHERE chain_id = $1;
