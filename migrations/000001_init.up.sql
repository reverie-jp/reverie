CREATE DOMAIN ulid AS TEXT CHECK (LENGTH(VALUE) = 26);

CREATE TABLE IF NOT EXISTS users (
    id ulid PRIMARY KEY,
    custom_id VARCHAR(15) NOT NULL UNIQUE,
    custom_id_change_time TIMESTAMPTZ,
    display_name VARCHAR(20) NOT NULL DEFAULT 'unknown',
    biography TEXT,
    location VARCHAR(100),
    website VARCHAR(255),
    avatar_url TEXT,
    banner_url TEXT,
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    following_count INT NOT NULL DEFAULT 0,
    follower_count INT NOT NULL DEFAULT 0,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
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

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id ulid PRIMARY KEY,
    user_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expire_time TIMESTAMPTZ NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

CREATE TYPE call_visibility AS ENUM ('open', 'users_only', 'locked');

CREATE TABLE IF NOT EXISTS calls (
    id ulid PRIMARY KEY,
    host_user_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    visibility call_visibility NOT NULL,
    end_time TIMESTAMPTZ,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_calls_host_user_id ON calls(host_user_id);
CREATE INDEX idx_calls_visibility_create_time ON calls(visibility, create_time DESC);

CREATE TABLE IF NOT EXISTS call_participants (
    call_id ulid NOT NULL REFERENCES calls(id) ON DELETE CASCADE,
    participant_identity TEXT NOT NULL,
    user_id ulid REFERENCES users(id) ON DELETE SET NULL,
    display_name VARCHAR(20) NOT NULL,
    first_join_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    disconnected_time TIMESTAMPTZ,
    muted_by_host BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (call_id, participant_identity)
);

CREATE INDEX idx_call_participants_order ON call_participants(call_id, first_join_time);
CREATE INDEX idx_call_participants_user_id ON call_participants(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_call_participants_active ON call_participants(call_id, last_seen_time) WHERE disconnected_time IS NULL;
CREATE INDEX idx_call_participants_active_user ON call_participants(user_id, last_seen_time) WHERE user_id IS NOT NULL AND disconnected_time IS NULL;

CREATE TABLE IF NOT EXISTS call_bans (
    call_id ulid NOT NULL REFERENCES calls(id) ON DELETE CASCADE,
    user_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (call_id, user_id)
);

CREATE TABLE IF NOT EXISTS user_follows (
    follower_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followee_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower_id, followee_id),
    CHECK (follower_id <> followee_id)
);

CREATE INDEX idx_user_follows_followee ON user_follows(followee_id, create_time DESC);
CREATE INDEX idx_user_follows_follower ON user_follows(follower_id, create_time DESC);

-- Maintain denormalized follow counts on users so reads don't scan user_follows.
CREATE FUNCTION user_follow_counts_sync() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE users SET following_count = following_count + 1 WHERE id = NEW.follower_id;
        UPDATE users SET follower_count = follower_count + 1 WHERE id = NEW.followee_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE users SET following_count = GREATEST(following_count - 1, 0) WHERE id = OLD.follower_id;
        UPDATE users SET follower_count = GREATEST(follower_count - 1, 0) WHERE id = OLD.followee_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_user_follows_counts
AFTER INSERT OR DELETE ON user_follows
FOR EACH ROW EXECUTE FUNCTION user_follow_counts_sync();
