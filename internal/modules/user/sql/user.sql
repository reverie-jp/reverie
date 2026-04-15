-- name: ListUsersByIDs :many
SELECT id, custom_id, custom_id_changed_at, display_name, biography,
       avatar_media_id, banner_media_id, is_private, birthdate, create_time, update_time
FROM users
WHERE id = ANY(@ids::text[]);

-- name: GetUserByCustomID :one
SELECT id, custom_id, custom_id_changed_at, display_name, biography,
       avatar_media_id, banner_media_id, is_private, birthdate, create_time, update_time
FROM users
WHERE custom_id = $1;

-- name: CreateUser :exec
INSERT INTO users (id, custom_id, display_name)
VALUES ($1, $2, $3);

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
