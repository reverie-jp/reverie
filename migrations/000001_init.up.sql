CREATE DOMAIN ulid AS TEXT CHECK (LENGTH(VALUE) = 26);

-- ================================================================
-- media
-- ================================================================

CREATE TYPE media_usage AS ENUM (
    'user_avatar',
    'user_banner',
    'post'
);

CREATE TYPE media_type AS ENUM (
    'image',
    'gif',
    'video',
    'audio'
);

CREATE TYPE media_status AS ENUM (
    'pending',
    'processing',
    'ready',
    'failed'
);

CREATE TABLE IF NOT EXISTS media (
    id                  ulid PRIMARY KEY,
    usage               media_usage NOT NULL,
    type                media_type DEFAULT NULL,
    status              media_status NOT NULL DEFAULT 'pending',
    width               INT DEFAULT NULL,
    height              INT DEFAULT NULL,
    original_image_url  TEXT DEFAULT NULL,
    thumbnail_image_url TEXT DEFAULT NULL,
    video_url           TEXT DEFAULT NULL,
    audio_url           TEXT DEFAULT NULL,
    duration_seconds    INT DEFAULT NULL,
    create_time         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ================================================================
-- users
-- ================================================================

CREATE TABLE IF NOT EXISTS users (
    id                   ulid PRIMARY KEY,
    custom_id            VARCHAR(15) NOT NULL UNIQUE,
    custom_id_changed_at TIMESTAMPTZ,
    display_name         VARCHAR(20) NOT NULL DEFAULT 'unknown',
    biography            TEXT NOT NULL DEFAULT '',
    avatar_media_id      ulid DEFAULT NULL REFERENCES media(id) ON DELETE SET NULL,
    banner_media_id      ulid DEFAULT NULL REFERENCES media(id) ON DELETE SET NULL,
    is_private           BOOLEAN NOT NULL DEFAULT FALSE,
    birthdate            TIMESTAMPTZ DEFAULT NULL,
    create_time          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TYPE auth_provider_type AS ENUM ('google', 'line');

CREATE TABLE IF NOT EXISTS auth_providers (
    id               ulid PRIMARY KEY,
    user_id          ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider         auth_provider_type NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    create_time      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (provider, provider_user_id)
);

CREATE INDEX idx_auth_providers_user_id ON auth_providers(user_id);

-- ================================================================
-- social graph
-- ================================================================

CREATE TABLE IF NOT EXISTS user_follows (
    follower_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followed_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower_id, followed_id)
);

CREATE INDEX idx_user_follows_followed_id ON user_follows(followed_id);

CREATE TABLE IF NOT EXISTS user_blocks (
    blocker_id  ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_id  ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (blocker_id, blocked_id)
);

CREATE INDEX idx_user_blocks_blocked_id ON user_blocks(blocked_id);

-- ================================================================
-- posts
-- ================================================================

CREATE TABLE IF NOT EXISTS posts (
    id          ulid PRIMARY KEY,
    author_id   ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reply_to_id ulid DEFAULT NULL REFERENCES posts(id) ON DELETE SET NULL,
    repost_id   ulid DEFAULT NULL REFERENCES posts(id) ON DELETE SET NULL,
    text        TEXT NOT NULL DEFAULT '',
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_posts_author_id  ON posts(author_id);
CREATE INDEX idx_posts_create_time ON posts(create_time DESC);

CREATE TABLE IF NOT EXISTS post_media (
    post_id  ulid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    media_id ulid NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    position INT NOT NULL CHECK (position BETWEEN 0 AND 3),
    PRIMARY KEY (post_id, media_id),
    UNIQUE (post_id, position)
);

CREATE TABLE IF NOT EXISTS post_favorites (
    user_id     ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id     ulid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);

CREATE INDEX idx_post_favorites_post_id ON post_favorites(post_id);

-- ================================================================
-- calls
-- ================================================================

CREATE TYPE call_type AS ENUM (
    'voice',
    'video'
);

CREATE TYPE call_joinable_by AS ENUM (
    'all',
    'followers',
    'friends',
    'nobody'
);

CREATE TYPE call_participant_role AS ENUM (
    'host',
    'co-host',
    'participant'
);

CREATE TABLE IF NOT EXISTS calls (
    id          ulid PRIMARY KEY,
    title       TEXT NOT NULL DEFAULT '',
    type        call_type NOT NULL DEFAULT 'voice',
    joinable_by call_joinable_by NOT NULL DEFAULT 'all',
    host_id     ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_time  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    end_time    TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_calls_host_id   ON calls(host_id);
CREATE INDEX idx_calls_end_time  ON calls(end_time) WHERE end_time IS NULL;

CREATE TABLE IF NOT EXISTS call_participants (
    call_id        ulid NOT NULL REFERENCES calls(id) ON DELETE CASCADE,
    participant_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role           call_participant_role NOT NULL DEFAULT 'participant',
    joined_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    left_at        TIMESTAMPTZ DEFAULT NULL,
    PRIMARY KEY (call_id, participant_id)
);

-- ================================================================
-- chat
-- ================================================================

CREATE TYPE room_type AS ENUM (
    'direct',
    'group'
);

CREATE TABLE IF NOT EXISTS rooms (
    id              ulid PRIMARY KEY,
    room_type       room_type NOT NULL,
    name            VARCHAR(100) DEFAULT NULL,
    group_image_url VARCHAR(255) DEFAULT NULL,
    create_time     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS room_members (
    id          ulid PRIMARY KEY,
    room_id     ulid NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    user_id     ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_read_at TIMESTAMPTZ,
    is_muted    BOOLEAN NOT NULL DEFAULT FALSE,
    mute_until  TIMESTAMPTZ DEFAULT NULL,
    UNIQUE (user_id, room_id)
);

CREATE INDEX idx_room_members_user_id ON room_members(user_id);
CREATE INDEX idx_room_members_room_id ON room_members(room_id);

CREATE TABLE IF NOT EXISTS pinned_rooms (
    id        ulid PRIMARY KEY,
    user_id   ulid NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    room_id   ulid NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    pinned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, room_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id          ulid PRIMARY KEY,
    room_id     ulid NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    sender_id   ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content     VARCHAR(10000),
    has_file    BOOLEAN NOT NULL DEFAULT FALSE,
    is_deleted  BOOLEAN NOT NULL DEFAULT FALSE,
    is_edited   BOOLEAN NOT NULL DEFAULT FALSE,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_room_id    ON messages(room_id);
CREATE INDEX idx_messages_create_time ON messages(room_id, create_time DESC);

CREATE TYPE file_type AS ENUM (
    'image',
    'video',
    'voice'
);

CREATE TABLE IF NOT EXISTS message_files (
    id          ulid PRIMARY KEY,
    message_id  ulid NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    media_id    ulid NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    file_type   file_type,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS message_reads (
    id         ulid PRIMARY KEY,
    message_id ulid NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    user_id    ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    read_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (message_id, user_id)
);

CREATE INDEX idx_message_reads_user_id ON message_reads(user_id);
