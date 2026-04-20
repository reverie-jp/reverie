-- drop tables
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS user_auth_providers;
DROP TABLE IF EXISTS users;

-- drop types
DROP TYPE IF EXISTS auth_provider;

-- drop domains
DROP DOMAIN IF EXISTS ulid;
