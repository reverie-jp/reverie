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

-- name: SearchUsers :many
SELECT id, custom_id, custom_id_changed_at, display_name, biography,
       avatar_media_id, banner_media_id, is_private, birthdate, create_time, update_time
FROM users
WHERE (custom_id ILIKE '%' || $1 || '%' OR display_name ILIKE '%' || $1 || '%')
  AND ($2::timestamptz IS NULL OR create_time < $2)
ORDER BY create_time DESC
LIMIT $3;

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

-- name: CreateUserFollow :exec
INSERT INTO user_follows (follower_id, followed_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeleteUserFollow :exec
DELETE FROM user_follows
WHERE follower_id = $1 AND followed_id = $2;

-- name: GetUserFollow :one
SELECT follower_id, followed_id, create_time
FROM user_follows
WHERE follower_id = $1 AND followed_id = $2;

-- name: CountUserFollowers :one
SELECT COUNT(*) FROM user_follows WHERE followed_id = $1;

-- name: CountUserFollowing :one
SELECT COUNT(*) FROM user_follows WHERE follower_id = $1;

-- name: ListFollowers :many
SELECT u.id, u.custom_id, u.custom_id_changed_at, u.display_name, u.biography,
       u.avatar_media_id, u.banner_media_id, u.is_private, u.birthdate, u.create_time, u.update_time
FROM users u
JOIN user_follows f ON f.follower_id = u.id
WHERE f.followed_id = $1
  AND ($2::timestamptz IS NULL OR f.create_time < $2)
ORDER BY f.create_time DESC
LIMIT $3;

-- name: ListFollowing :many
SELECT u.id, u.custom_id, u.custom_id_changed_at, u.display_name, u.biography,
       u.avatar_media_id, u.banner_media_id, u.is_private, u.birthdate, u.create_time, u.update_time
FROM users u
JOIN user_follows f ON f.followed_id = u.id
WHERE f.follower_id = $1
  AND ($2::timestamptz IS NULL OR f.create_time < $2)
ORDER BY f.create_time DESC
LIMIT $3;
