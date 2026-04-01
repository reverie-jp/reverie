-- name: CreateUser :exec
INSERT INTO users (id, custom_id, display_name, avatar_url, create_time, update_time)
VALUES ($1, $2, $3, $4, $5, $5);

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND delete_time IS NULL;

-- name: GetUserByCustomID :one
SELECT * FROM users
WHERE custom_id = $1 AND delete_time IS NULL;

-- name: SoftDeleteUser :exec
UPDATE users SET delete_time = NOW() WHERE id = $1;

-- name: CreateAuthProvider :exec
INSERT INTO user_auth_providers (id, user_id, provider, provider_user_id, create_time)
VALUES ($1, $2, $3, $4, $5);

-- name: GetAuthProviderByProvider :one
SELECT * FROM user_auth_providers
WHERE provider = $1 AND provider_user_id = $2;
