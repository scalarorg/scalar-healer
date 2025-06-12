CREATE TABLE IF NOT EXISTS bridge_requests (
    id BIGSERIAL PRIMARY KEY,
    address BYTEA NOT NULL,
    signature BYTEA NOT NULL,
    chain TEXT NOT NULL,
    tx_hash BYTEA NOT NULL,
    nonce NUMERIC NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS redeem_requests (
    id BIGSERIAL PRIMARY KEY,
    address BYTEA NOT NULL,
    source_chain TEXT NOT NULL,
    dest_chain TEXT NOT NULL,
    symbol TEXT NOT NULL,
    amount TEXT NOT NULL,
    locking_script BYTEA NOT NULL,
    custodian_group_uid BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transfer_requests (
    id BIGSERIAL PRIMARY KEY,
    address BYTEA NOT NULL,
    signature BYTEA NOT NULL,
    chain TEXT NOT NULL,
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
    chain TEXT NOT NULL,
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
    custodians JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS protocols (
    id BIGSERIAL PRIMARY KEY,

    symbol TEXT UNIQUE NOT NULL CHECK (symbol <> ''),
    name TEXT NOT NULL,
    bitcoin_pubkey BYTEA NOT NULL,
    decimals BIGINT NOT NULL,
    avatar TEXT NOT NULL,
    
    custodian_group_name TEXT NOT NULL,
    custodian_group_uid BYTEA NOT NULL,
    tag TEXT NOT NULL,
    
    liquidity_model TEXT NOT NULL,
   
    capacity NUMERIC NOT NULL,
    daily_mint_limit NUMERIC NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tokens (
    id BIGSERIAL PRIMARY KEY,

    symbol TEXT NOT NULL,
    chain TEXT NOT NULL,
    chain_id NUMERIC NOT NULL,
    active BOOLEAN NOT NULL,
    address BYTEA NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS tokens_symbol_chain_idx ON tokens (symbol, chain);

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
    request_id BYTEA PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS utxo_reservations (
    utxo_tx_id BYTEA NOT NULL,
    utxo_vout BIGINT NOT NULL,
    reservation_id BYTEA NOT NULL,
    amount NUMERIC NOT NULL,

    PRIMARY KEY (utxo_tx_id, utxo_vout, reservation_id, amount),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Redeem session

CREATE TYPE REDEEM_PHASE as ENUM ('PREPARING', 'EXECUTING'); 

CREATE TABLE IF NOT EXISTS redeem_sessions (
    id BIGSERIAL PRIMARY KEY,
    custodian_group_uid BYTEA UNIQUE NOT NULL,
    sequence BIGINT NOT NULL CHECK (sequence >= 0),
    current_phase REDEEM_PHASE NOT NULL,
    last_redeem_tx BYTEA,
    is_switching BOOLEAN,
    phase_expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chain_redeem_sessions (
    id BIGSERIAL PRIMARY KEY,
    
    chain TEXT NOT NULL,
    custodian_group_uid BYTEA NOT NULL,
    sequence BIGINT NOT NULL CHECK (sequence >= 0),
    current_phase REDEEM_PHASE NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS chain_redeem_sessions_custodian_group_uid_chain_idx ON chain_redeem_sessions (custodian_group_uid, chain);

-- Commands

CREATE TYPE COMMAND_TYPE as ENUM (
   'burnToken',
   'deployToken',
   'mintToken',
   'approveContractCall',
   'approveContractCallWithMint',
   'transferOperatorship',
   'switchPhase',
   'registerCustodianGroup',
   'redeemToken'
);

CREATE TABLE IF NOT EXISTS commands (
    id BYTEA PRIMARY KEY,
    chain TEXT NOT NULL,
    command_batch_id BYTEA NULL,
    payload BYTEA NOT NULL,
    params BYTEA NOT NULL,
    status INT CHECK (status IN (0, 1, 2)),
    -- 0 = PENDING, not included in the batch command
    -- 1 = QUEUED, included in the batch command
    -- 2 = EXECUTED, executed
    command_type COMMAND_TYPE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS commands_chain ON commands (chain);

CREATE TYPE BATCH_STATUS as ENUM ('PENDING', 'SIGNED', 'EXECUTED'); 

CREATE TABLE IF NOT EXISTS command_batchs (
    id BYTEA PRIMARY KEY,
    chain TEXT NOT NULL,
    data BYTEA NOT NULL,
    sig_hash BYTEA NOT NULL,
    signature BYTEA,
    status BATCH_STATUS NOT NULL,
    extra_data BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS command_batchs_chain ON command_batchs (chain);

-- Redeem Command is special command that is not included in the batch command, it is a single command but we have to handle it as a batch command

CREATE TABLE IF NOT EXISTS redeem_commands (
    id BYTEA PRIMARY KEY,
    request_id BIGINT NOT NULL,
    chain TEXT NOT NULL,
    status BATCH_STATUS NOT NULL,
    params BYTEA NOT NULL,
    data BYTEA NOT NULL,
    sig_hash BYTEA NOT NULL,
    signature BYTEA,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS redeem_commands_chain ON redeem_commands (chain);

ALTER TABLE utxo_reservations ADD FOREIGN KEY (reservation_id) REFERENCES reservations(request_id) ON DELETE CASCADE;

ALTER TABLE utxo_reservations ADD FOREIGN KEY (utxo_tx_id, utxo_vout) REFERENCES utxos(tx_id, vout) ON DELETE CASCADE;

ALTER TABLE tokens ADD FOREIGN KEY (symbol) REFERENCES protocols (symbol) ON DELETE CASCADE;

ALTER TABLE protocols ADD FOREIGN KEY (custodian_group_uid) REFERENCES custodian_groups (uid) ON DELETE CASCADE;

ALTER TABLE redeem_sessions ADD FOREIGN KEY (custodian_group_uid) REFERENCES custodian_groups (uid) ON DELETE CASCADE;

ALTER TABLE chain_redeem_sessions ADD FOREIGN KEY (custodian_group_uid) REFERENCES custodian_groups (uid) ON DELETE CASCADE;

ALTER TABLE redeem_requests ADD FOREIGN KEY (custodian_group_uid) REFERENCES custodian_groups (uid) ON DELETE CASCADE;

ALTER TABLE commands ADD FOREIGN KEY (command_batch_id) REFERENCES command_batchs (id);

ALTER TABLE redeem_commands ADD FOREIGN KEY (request_id) REFERENCES redeem_requests (id) ON DELETE CASCADE;