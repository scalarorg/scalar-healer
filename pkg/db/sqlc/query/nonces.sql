-- name: UpsertNonce :exec
INSERT INTO nonces (address,nonce) VALUES ($1,$2) 
ON CONFLICT (address) DO UPDATE SET nonce = EXCLUDED.nonce;

-- name: GetNonce :one
SELECT nonce FROM nonces
WHERE address = $1;