-- name: CreateNotification :one
-- Inserts a notification if the dedup key is free. Returns the existing row
-- on conflict, so callers always get the canonical record to publish downstream.
INSERT INTO notifications (id, recipient_user_id, type, actor_user_id, resource_name)
VALUES (
    sqlc.arg(id)::ulid,
    sqlc.arg(recipient_user_id)::ulid,
    sqlc.arg(type)::notification_type,
    sqlc.narg(actor_user_id)::ulid,
    sqlc.arg(resource_name)::text
)
ON CONFLICT (recipient_user_id, type, COALESCE(actor_user_id::text, ''), resource_name)
DO UPDATE SET id = notifications.id
RETURNING *;

-- name: ListNotificationsByRecipient :many
-- Keyset pagination by id DESC (ULID).
SELECT * FROM notifications
WHERE recipient_user_id = sqlc.arg(recipient_user_id)::ulid
  AND (sqlc.arg(cursor_id)::text = '' OR id < sqlc.arg(cursor_id)::ulid)
ORDER BY id DESC
LIMIT sqlc.arg(page_size)::int;

-- name: MarkNotificationsRead :execrows
UPDATE notifications
SET read_time = NOW()
WHERE recipient_user_id = sqlc.arg(recipient_user_id)::ulid
  AND read_time IS NULL
  AND id = ANY(sqlc.arg(ids)::text[]);

-- name: MarkAllNotificationsRead :execrows
UPDATE notifications
SET read_time = NOW()
WHERE recipient_user_id = sqlc.arg(recipient_user_id)::ulid
  AND read_time IS NULL;

-- name: CountUnreadNotifications :one
SELECT COUNT(*)::int AS count
FROM notifications
WHERE recipient_user_id = sqlc.arg(recipient_user_id)::ulid
  AND read_time IS NULL;

-- name: GetNotification :one
SELECT * FROM notifications
WHERE id = sqlc.arg(id)::ulid;

-- name: CreateFanOutNotifications :many
-- Batched insert for "same actor/type/resource, many recipients" fan-out
-- (e.g. call_started broadcast to followers). Row IDs are client-generated
-- ULIDs passed via parallel arrays. ON CONFLICT DO NOTHING: conflicted rows
-- (already notified) are silently skipped and omitted from RETURNING, so
-- callers only publish events for genuinely new notifications.
INSERT INTO notifications (id, recipient_user_id, type, actor_user_id, resource_name)
SELECT
    x.id::ulid,
    x.recipient_id::ulid,
    sqlc.arg(type)::notification_type,
    sqlc.narg(actor_user_id)::ulid,
    sqlc.arg(resource_name)::text
FROM unnest(sqlc.arg(ids)::text[], sqlc.arg(recipient_ids)::text[]) AS x(id, recipient_id)
ON CONFLICT (recipient_user_id, type, COALESCE(actor_user_id::text, ''), resource_name)
DO NOTHING
RETURNING *;

-- name: DeleteNotificationsByTypeActor :execrows
-- Removes notifications matching (recipient, type, actor). Used to clear
-- user_followed notifications when the follow relationship is undone, so a
-- re-follow creates a fresh notification instead of hitting the dedup index.
DELETE FROM notifications
WHERE recipient_user_id = sqlc.arg(recipient_user_id)::ulid
  AND type = sqlc.arg(type)::notification_type
  AND actor_user_id = sqlc.arg(actor_user_id)::ulid;
