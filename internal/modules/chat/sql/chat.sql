-- name: ListRoomsByUser :many
SELECT
  r.id,
  r.room_type,
  r.name,
  rm.is_muted,
  rm.last_read_at,
  EXISTS(SELECT 1 FROM pinned_rooms pr WHERE pr.user_id = $1 AND pr.room_id = r.id) AS is_pinned,
  (SELECT m.content FROM messages m WHERE m.room_id = r.id AND NOT m.is_deleted ORDER BY m.create_time DESC LIMIT 1) AS last_message_text,
  (SELECT m.sender_id FROM messages m WHERE m.room_id = r.id AND NOT m.is_deleted ORDER BY m.create_time DESC LIMIT 1) AS last_message_sender_id,
  (SELECT m.create_time FROM messages m WHERE m.room_id = r.id AND NOT m.is_deleted ORDER BY m.create_time DESC LIMIT 1) AS last_message_at,
  (SELECT COUNT(*) FROM messages m WHERE m.room_id = r.id AND NOT m.is_deleted AND m.create_time > COALESCE(rm.last_read_at, '-infinity'::timestamptz)) AS unread_count
FROM rooms r
JOIN room_members rm ON rm.room_id = r.id AND rm.user_id = $1
ORDER BY last_message_at DESC NULLS LAST;

-- name: GetRoomByID :one
SELECT * FROM rooms WHERE id = $1;

-- name: GetDirectRoom :one
SELECT r.id FROM rooms r
JOIN room_members rm1 ON rm1.room_id = r.id AND rm1.user_id = $1
JOIN room_members rm2 ON rm2.room_id = r.id AND rm2.user_id = $2
WHERE r.room_type = 'direct'
LIMIT 1;

-- name: CreateRoom :one
INSERT INTO rooms (id, room_type, name) VALUES ($1, $2::room_type, $3) RETURNING *;

-- name: CreateRoomMember :exec
INSERT INTO room_members (id, room_id, user_id) VALUES ($1, $2, $3);

-- name: GetRoomMember :one
SELECT * FROM room_members WHERE room_id = $1 AND user_id = $2;

-- name: GetRoomOtherMember :one
SELECT user_id FROM room_members WHERE room_id = $1 AND user_id != $2 LIMIT 1;

-- name: ListMessages :many
SELECT * FROM messages
WHERE room_id = $1 AND NOT is_deleted
  AND ($2::timestamptz IS NULL OR create_time < $2)
ORDER BY create_time DESC
LIMIT $3;

-- name: CreateMessage :one
INSERT INTO messages (id, room_id, sender_id, content) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: MarkRoomAsRead :exec
UPDATE room_members SET last_read_at = NOW() WHERE room_id = $1 AND user_id = $2;

-- name: PinRoom :exec
INSERT INTO pinned_rooms (id, user_id, room_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;

-- name: UnpinRoom :exec
DELETE FROM pinned_rooms WHERE user_id = $1 AND room_id = $2;

-- name: MuteRoom :exec
UPDATE room_members SET is_muted = $3 WHERE room_id = $1 AND user_id = $2;

-- name: LeaveRoom :exec
DELETE FROM room_members WHERE room_id = $1 AND user_id = $2;

-- name: GetRoomMembers :many
SELECT user_id FROM room_members WHERE room_id = $1;

-- name: RemoveRoomMember :exec
DELETE FROM room_members WHERE room_id = $1 AND user_id = $2;

-- name: UpdateRoomName :one
UPDATE rooms SET name = $2, update_time = NOW() WHERE id = $1 RETURNING id, room_type, name, group_image_url, create_time, update_time;
