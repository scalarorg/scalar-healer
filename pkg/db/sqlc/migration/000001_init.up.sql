CREATE TABLE IF NOT EXISTS bridge_requests (
    id BIGSERIAL PRIMARY KEY,
    address BYTEA NOT NULL,
    signature BYTEA NOT NULL,
    chain_id NUMERIC NOT NULL,
    tx_hash BYTEA NOT NULL,
    nonce NUMERIC NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS redeem_requests (
    id BIGSERIAL PRIMARY KEY,
    address BYTEA NOT NULL,
    signature BYTEA NOT NULL,
    chain_id NUMERIC NOT NULL,
    symbol TEXT NOT NULL,
    amount TEXT NOT NULL,
    nonce NUMERIC NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transfer_requests (
    id BIGSERIAL PRIMARY KEY,
    address BYTEA NOT NULL,
    signature BYTEA NOT NULL,
    chain_id NUMERIC NOT NULL,
    destination_chain TEXT NOT NULL,
    destination_address BYTEA NOT NULL,
    symbol TEXT NOT NULL,
    amount TEXT NOT NULL,
    nonce NUMERIC NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gateway_addresses (
    id BIGSERIAL PRIMARY KEY,
    address BYTEA NOT NULL,
    chain_id NUMERIC NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tokens (
    id BIGSERIAL PRIMARY KEY,
    protocol TEXT NOT NULL,
    symbol TEXT NOT NULL,
    chain_id NUMERIC NOT NULL,
    active BOOLEAN NOT NULL,
    address BYTEA NOT NULL,
    decimal NUMERIC NOT NULL,
    name TEXT NOT NULL,
    avatar TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS nonces (
    id BIGSERIAL PRIMARY KEY,
    address BYTEA UNIQUE NOT NULL,
    nonce NUMERIC NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS nonces_address_idx ON nonces (address);

CREATE TABLE IF NOT EXISTS custodian_groups (
    id BIGSERIAL PRIMARY KEY,
    uid BYTEA UNIQUE NOT NULL,
    name TEXT NOT NULL,
    bitcoin_pubkey BYTEA NOT NULL,
    quorum BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS protocols (
    id BIGSERIAL PRIMARY KEY,
    asset TEXT UNIQUE NOT NULL CHECK (asset <> ''),
    name TEXT NOT NULL,
    custodian_group_name TEXT NOT NULL,
    custodian_group_uid BYTEA UNIQUE NOT NULL,
    tag TEXT NOT NULL,
    liquidity_model TEXT NOT NULL,
    symbol TEXT NOT NULL,
    decimals BIGINT NOT NULL,
    capacity NUMERIC NOT NULL,
    daily_mint_limit NUMERIC NOT NULL,
    avatar TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS protocols_asset_idx ON protocols (asset);