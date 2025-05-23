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
    symbol TEXT NOT NULL,
    name TEXT NOT NULL,
    chain_id NUMERIC NOT NULL,
    active BOOLEAN NOT NULL,
    address BYTEA NOT NULL,
    decimal NUMERIC NOT NULL,
    avatar TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS tokens_symbol_chain_id_idx ON tokens (symbol, chain_id);

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
    symbol TEXT UNIQUE NOT NULL CHECK (symbol <> ''),
    name TEXT NOT NULL,
    custodian_group_name TEXT NOT NULL,
    custodian_group_uid BYTEA UNIQUE NOT NULL,
    tag TEXT NOT NULL,
    liquidity_model TEXT NOT NULL,
    decimals BIGINT NOT NULL,
    capacity NUMERIC NOT NULL,
    daily_mint_limit NUMERIC NOT NULL,
    avatar TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS protocols_symbol_idx ON protocols (symbol);

CREATE TABLE IF NOT EXISTS utxos (
    id BIGSERIAL PRIMARY KEY,
    tx_id BYTEA NOT NULL,
    vout BIGINT NOT NULL,
    script_pubkey BYTEA NOT NULL,
    amount_in_sats NUMERIC NOT NULL,
    custodian_group_uid BYTEA NOT NULL,
    block_height BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS utxos_tx_id_vout_idx ON utxos (tx_id, vout);

CREATE INDEX IF NOT EXISTS utxos_custodian_group_uid_idx ON utxos (custodian_group_uid);

CREATE TABLE IF NOT EXISTS reservations (
    id BIGSERIAL PRIMARY KEY,
    utxo_tx_id BYTEA NOT NULL,
    utxo_vout BIGINT NOT NULL,
    request_id TEXT UNIQUE NOT NULL,
    amount NUMERIC NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS reservations_utxo_tx_id_vout_idx ON reservations (utxo_tx_id, utxo_vout);

ALTER TABLE reservations ADD FOREIGN KEY (utxo_tx_id, utxo_vout) REFERENCES utxos (tx_id, vout) ON DELETE CASCADE;


-- Redeem session

CREATE TYPE REDEEM_PHASE as ENUM ('PREPARING', 'EXECUTING'); 

CREATE TABLE IF NOT EXISTS redeem_sessions (
    id BIGSERIAL PRIMARY KEY,
    custodian_group_uid BYTEA NOT NULL,
    sequence BIGINT NOT NULL CHECK (sequence >= 0),
    chain TEXT NOT NULL,
    current_phase REDEEM_PHASE NOT NULL,
    last_redeem_tx BYTEA NOT NULL,
    is_switching BOOLEAN NOT NULL,
    phase_expired_at BIGINT NOT NULL CHECK (phase_expired_at >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE redeem_sessions ADD FOREIGN KEY (custodian_group_uid) REFERENCES custodian_groups (uid) ON DELETE CASCADE;
CREATE UNIQUE INDEX IF NOT EXISTS redeem_sessions_custodian_group_uid_chain_idx ON redeem_sessions (custodian_group_uid, chain);
