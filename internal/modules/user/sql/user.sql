-- name: ListUsersByIDs :many
SELECT id, custom_id, custom_id_changed_at, display_name, biography,
       avatar_media_id, banner_media_id, is_private, birthdate, create_time, update_time
FROM users
WHERE id = ANY(@ids::text[]);

-- name: GetUserByID :one
SELECT id, custom_id, custom_id_changed_at, display_name, biography,
       avatar_media_id, banner_media_id, is_private, birthdate, create_time, update_time
FROM users
WHERE id = $1;

-- name: GetUserByCustomID :one
SELECT id, custom_id, custom_id_changed_at, display_name, biography,
       avatar_media_id, banner_media_id, is_private, birthdate, create_time, update_time
FROM users
WHERE custom_id = $1;

-- name: UpdateUser :one
UPDATE users
SET
  display_name = $2,
  biography    = $3,
  is_private   = $4,
  birthdate    = $5,
  update_time  = NOW()
WHERE id = $1
RETURNING id, custom_id, custom_id_changed_at, display_name, biography,
          avatar_media_id, banner_media_id, is_private, birthdate, create_time, update_time;

-- name: CreateUser :exec
INSERT INTO users (id, custom_id, display_name)
VALUES ($1, $2, $3);

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
