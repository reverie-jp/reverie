-- name: ListUsersByIDs :many
SELECT * FROM users
WHERE id = ANY(@ids::text[]);

-- name: GetUserByCustomID :one
SELECT * FROM users
WHERE custom_id = $1;

-- name: CreateUser :exec
INSERT INTO users (id, custom_id, display_name, avatar_url)
VALUES ($1, $2, $3, $4);

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUserLastSeen :exec
-- Called by the PresenceService.Heartbeat RPC every ~30s. `last_seen_time`
-- is not indexed, so Postgres can apply this as a HOT update and avoid
-- touching users' other indexes.
UPDATE users SET last_seen_time = NOW() WHERE id = sqlc.arg(id)::ulid;
