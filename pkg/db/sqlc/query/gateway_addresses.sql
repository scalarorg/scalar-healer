-- name: CreateGatewayAddress :exec
INSERT INTO gateway_addresses (address, chain)
VALUES ($1, $2);

-- name: CreateGatewayAddresses :exec
INSERT INTO gateway_addresses (address, chain)
VALUES (unnest($1::bytea[]), unnest($2::text[]));

-- name: GetGatewayAddress :one
SELECT address
FROM gateway_addresses
WHERE chain = $1;
