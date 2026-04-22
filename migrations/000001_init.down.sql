-- drop triggers and functions
DROP TRIGGER IF EXISTS trg_user_follows_counts ON user_follows;
DROP FUNCTION IF EXISTS user_follow_counts_sync();

-- drop tables
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS user_follows;
DROP TABLE IF EXISTS call_bans;
DROP TABLE IF EXISTS call_participants;
DROP TABLE IF EXISTS calls;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS user_auth_providers;
DROP TABLE IF EXISTS users;

-- drop types
DROP TYPE IF EXISTS notification_type;
DROP TYPE IF EXISTS call_visibility;
DROP TYPE IF EXISTS auth_provider;

-- drop domains
DROP DOMAIN IF EXISTS ulid;
