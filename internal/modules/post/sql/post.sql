-- name: CreatePost :one
INSERT INTO posts (id, author_id, short_id, reply_to_post_id, repost_post_id, text)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, author_id, short_id, reply_to_post_id, repost_post_id, text, create_time, update_time;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1 AND author_id = $2;

-- name: GetPostByID :one
SELECT id, author_id, short_id, reply_to_post_id, repost_post_id, text, create_time, update_time
FROM posts
WHERE id = $1;

-- name: GetPostByShortID :one
SELECT p.id, p.author_id, p.short_id, p.reply_to_post_id, p.repost_post_id, p.text, p.create_time, p.update_time
FROM posts p
JOIN users u ON u.id = p.author_id
WHERE u.custom_id = $1
  AND p.short_id = $2;

-- name: ListTimeline :many
SELECT id, author_id, short_id, reply_to_post_id, repost_post_id, text, create_time, update_time
FROM posts
WHERE reply_to_post_id IS NULL
  AND ($1::timestamptz IS NULL OR create_time < $1)
ORDER BY create_time DESC
LIMIT $2;

-- name: ListFollowingTimeline :many
SELECT id, author_id, short_id, reply_to_post_id, repost_post_id, text, create_time, update_time
FROM posts
WHERE reply_to_post_id IS NULL
  AND (
    author_id IN (SELECT followee_id FROM user_follows WHERE follower_id = $1)
    OR author_id = $1
  )
  AND ($2::timestamptz IS NULL OR create_time < $2)
ORDER BY create_time DESC
LIMIT $3;

-- name: ListUserPosts :many
SELECT id, author_id, short_id, reply_to_post_id, repost_post_id, text, create_time, update_time
FROM posts
WHERE author_id = $1
  AND reply_to_post_id IS NULL
  AND ($2::timestamptz IS NULL OR create_time < $2)
ORDER BY create_time DESC
LIMIT $3;

-- name: ListPostReplies :many
SELECT id, author_id, short_id, reply_to_post_id, repost_post_id, text, create_time, update_time
FROM posts
WHERE reply_to_post_id = $1
  AND ($2::timestamptz IS NULL OR create_time < $2)
ORDER BY create_time DESC
LIMIT $3;

-- name: CountPostReplies :one
SELECT COUNT(*) FROM posts WHERE reply_to_post_id = $1;

-- name: ListPostReposts :many
SELECT id, author_id, short_id, reply_to_post_id, repost_post_id, text, create_time, update_time
FROM posts
WHERE repost_post_id = $1
  AND ($2::timestamptz IS NULL OR create_time < $2)
ORDER BY create_time DESC
LIMIT $3;

-- name: CountPostReposts :one
SELECT COUNT(*) FROM posts WHERE repost_post_id = $1;

-- name: CountPostFavorites :one
SELECT COUNT(*) FROM post_favorites WHERE post_id = $1;

-- name: GetPostFavorite :one
SELECT user_id, post_id FROM post_favorites
WHERE user_id = $1 AND post_id = $2;

-- name: CreatePostFavorite :exec
INSERT INTO post_favorites (user_id, post_id) VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeletePostFavorite :exec
DELETE FROM post_favorites WHERE user_id = $1 AND post_id = $2;

-- name: ListUserLikedPosts :many
SELECT p.id, p.author_id, p.short_id, p.reply_to_post_id, p.repost_post_id, p.text, p.create_time, p.update_time
FROM post_favorites pf
JOIN posts p ON p.id = pf.post_id
WHERE pf.user_id = $1
  AND ($2::timestamptz IS NULL OR pf.create_time < $2)
ORDER BY pf.create_time DESC
LIMIT $3;

-- name: ListPostLikes :many
SELECT u.id, u.custom_id, u.custom_id_change_time, u.display_name, u.biography, u.location, u.website, u.avatar_url, u.banner_url, u.is_private, u.following_count, u.follower_count, u.last_seen_time, u.create_time, u.update_time
FROM post_favorites pf
JOIN users u ON u.id = pf.user_id
WHERE pf.post_id = $1
ORDER BY pf.create_time DESC
LIMIT $2;
