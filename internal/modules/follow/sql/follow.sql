-- name: CreateUserFollow :exec
INSERT INTO user_follows (follower_id, followee_id)
VALUES (sqlc.arg(follower_id)::ulid, sqlc.arg(followee_id)::ulid)
ON CONFLICT DO NOTHING;

-- name: DeleteUserFollow :exec
DELETE FROM user_follows
WHERE follower_id = sqlc.arg(follower_id)::ulid
  AND followee_id = sqlc.arg(followee_id)::ulid;

-- name: IsFollowing :one
SELECT EXISTS (
  SELECT 1 FROM user_follows
  WHERE follower_id = sqlc.arg(follower_id)::ulid
    AND followee_id = sqlc.arg(followee_id)::ulid
) AS following;

-- name: ListFollowingIDs :many
-- Keyset pagination by followee_id DESC (ULID). cursor_id="" => first page.
SELECT followee_id FROM user_follows
WHERE follower_id = sqlc.arg(follower_id)::ulid
  AND (sqlc.arg(cursor_id)::text = '' OR followee_id < sqlc.arg(cursor_id)::ulid)
ORDER BY followee_id DESC
LIMIT sqlc.arg(page_size)::int;

-- name: ListFollowerIDs :many
SELECT follower_id FROM user_follows
WHERE followee_id = sqlc.arg(followee_id)::ulid
  AND (sqlc.arg(cursor_id)::text = '' OR follower_id < sqlc.arg(cursor_id)::ulid)
ORDER BY follower_id DESC
LIMIT sqlc.arg(page_size)::int;

-- name: ListAllFollowerIDs :many
-- Returns every follower_id for the given followee in one shot. Used by
-- fan-out writers (e.g. call creation) that need to touch all followers
-- rather than paginate. At reverie's scale this is cheaper than multiple
-- paginated round-trips.
SELECT follower_id FROM user_follows
WHERE followee_id = sqlc.arg(followee_id)::ulid;

-- name: ListFollowingEdgesForRequester :many
-- IDs from target set that the requester follows. Used for batched
-- is_following computation.
SELECT followee_id FROM user_follows
WHERE follower_id = sqlc.arg(follower_id)::ulid
  AND followee_id = ANY(sqlc.arg(followee_ids)::text[]);

-- name: ListFollowerEdgesForRequester :many
-- IDs from target set that follow the requester. Used for batched
-- is_followed_by computation.
SELECT follower_id FROM user_follows
WHERE followee_id = sqlc.arg(followee_id)::ulid
  AND follower_id = ANY(sqlc.arg(follower_ids)::text[]);

