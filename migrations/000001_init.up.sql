CREATE DOMAIN ulid AS TEXT CHECK (LENGTH(VALUE) = 26);

CREATE TABLE IF NOT EXISTS users (
    id ulid PRIMARY KEY,
    custom_id VARCHAR(15) NOT NULL UNIQUE,
    custom_id_changed_at TIMESTAMPTZ,
    display_name VARCHAR(20) NOT NULL DEFAULT 'unknown',
    biography TEXT,
    location VARCHAR(100),
    website VARCHAR(255),
    avatar_url TEXT,
    banner_url TEXT,
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    delete_time TIMESTAMPTZ
);

CREATE TYPE auth_provider AS ENUM ('google');

CREATE TABLE IF NOT EXISTS user_auth_providers (
    id ulid PRIMARY KEY,
    user_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider auth_provider NOT NULL,
    provider_user_id TEXT NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_user_auth_providers_user_id ON user_auth_providers(user_id);
