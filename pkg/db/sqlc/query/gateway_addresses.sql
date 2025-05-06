-- name: CreateGatewayAddress :exec
INSERT INTO gateway_addresses (address, chain_id)
VALUES ($1, $2);

-- name: GetGatewayAddress :one
SELECT address
FROM gateway_addresses
WHERE chain_id = $1;
