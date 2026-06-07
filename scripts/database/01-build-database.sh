#!/usr/bin/env bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL

    CREATE TABLE IF NOT EXISTS users (
        user_id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
        email                 TEXT        NOT NULL UNIQUE,
        public_key            BYTEA       NOT NULL,
        public_key_created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        last_login_at         TIMESTAMPTZ,
        deleted_at            TIMESTAMPTZ,
        created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at            TIMESTAMPTZ NOT NULL DEFAULT now()
    );

    CREATE TABLE IF NOT EXISTS sessions (
        session_id  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id     UUID        NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
        expires_at  TIMESTAMPTZ NOT NULL,
        revoked_at  TIMESTAMPTZ,
        ip_address  INET,
        user_agent  TEXT,
        created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
    );
    CREATE TABLE IF NOT EXISTS vault (
        vault_id      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
        vault_name    TEXT        NOT NULL,
        created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
        deleted_at    TIMESTAMPTZ,
        user_id       UUID        NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT
    );
    
    CREATE TABLE IF NOT EXISTS vault_entry (
        entry_id      UUID        PRIMARY KEY,
        entry_blob    BYTEA       NOT NULL,
        created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
        deleted_at    TIMESTAMPTZ,
        vault_id      UUID        NOT NULL REFERENCES vault(vault_id) ON DELETE RESTRICT

    );
    CREATE TABLE IF NOT EXISTS audit_log_entry (
        log_id        UUID        NOT NULL DEFAULT gen_random_uuid(),
        log_timestamp TIMESTAMPTZ NOT NULL,
        actor_id      UUID,
        actor_type    TEXT        NOT NULL,
        action        TEXT        NOT NULL,
        resource      TEXT,
        resource_id   TEXT,
        outcome       TEXT        NOT NULL,
        ip_address    INET,
        metadata      JSONB,
        PRIMARY KEY (log_id, log_timestamp)
    );
    SELECT create_hypertable(
        'audit_log_entry',
        'log_timestamp',
        chunk_time_interval => INTERVAL '7 days',
        if_not_exists => TRUE
    );

    ALTER TABLE audit_log_entry SET (
        timescaledb.compress,
        timescaledb.compress_orderby = 'log_timestamp DESC'
    );

    SELECT add_compression_policy('audit_log_entry', INTERVAL '14 days');
    
    CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
    CREATE INDEX IF NOT EXISTS idx_vault_user_id ON vault(user_id);

EOSQL