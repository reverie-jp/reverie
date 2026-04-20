CREATE TABLE IF NOT EXISTS message_reactions (
    id         ulid PRIMARY KEY,
    message_id ulid NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    user_id    ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    emoji      VARCHAR(10) NOT NULL DEFAULT '❤️',
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (message_id, user_id, emoji)
);

CREATE INDEX idx_message_reactions_message_id ON message_reactions(message_id);
