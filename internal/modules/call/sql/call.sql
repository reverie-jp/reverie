-- name: CreateCall :exec
INSERT INTO calls (id, host_user_id, visibility)
VALUES ($1, $2, $3);

-- name: GetCall :one
SELECT * FROM calls
WHERE id = sqlc.arg(id)::ulid;

-- name: UpdateCallVisibility :exec
UPDATE calls
SET visibility = sqlc.arg(visibility), update_time = NOW()
WHERE id = sqlc.arg(id)::ulid;

-- name: ListActivePublicCalls :many
-- Returns all active non-hidden calls (OPEN and USERS_ONLY). The usecase
-- filters further based on the caller's auth state. Keyset paginated by
-- ULID (monotonic, DESC). cursor_id="" means first page.
SELECT DISTINCT c.* FROM calls c
WHERE c.visibility IN ('open', 'users_only')
  AND EXISTS (
    SELECT 1 FROM call_participants p
    WHERE p.call_id = c.id
      AND p.last_seen_time > NOW() - (sqlc.arg(stale_seconds)::int || ' seconds')::interval
      AND p.disconnected_time IS NULL
  )
  AND (sqlc.arg(cursor_id)::text = '' OR c.id < sqlc.arg(cursor_id)::ulid)
ORDER BY c.id DESC
LIMIT sqlc.arg(page_size)::int;

-- name: UpsertCallParticipant :exec
INSERT INTO call_participants (call_id, participant_identity, user_id, display_name, last_seen_time, disconnected_time)
VALUES ($1, $2, $3, $4, NOW(), NULL)
ON CONFLICT (call_id, participant_identity) DO UPDATE
SET last_seen_time = NOW(),
    disconnected_time = NULL,
    display_name = EXCLUDED.display_name;

-- name: HeartbeatCallParticipant :execrows
UPDATE call_participants
SET last_seen_time = NOW()
WHERE call_id = sqlc.arg(call_id)::ulid
  AND participant_identity = sqlc.arg(participant_identity)
  AND disconnected_time IS NULL;

-- name: MarkCallParticipantDisconnected :execrows
UPDATE call_participants
SET disconnected_time = NOW()
WHERE call_id = sqlc.arg(call_id)::ulid
  AND participant_identity = sqlc.arg(participant_identity)
  AND disconnected_time IS NULL;

-- name: ListCallParticipants :many
SELECT * FROM call_participants
WHERE call_id = sqlc.arg(call_id)::ulid
ORDER BY first_join_time ASC;

-- name: GetActiveCallByUser :one
SELECT c.* FROM calls c
JOIN call_participants p ON p.call_id = c.id
WHERE p.user_id = sqlc.arg(user_id)::ulid
  AND p.last_seen_time > NOW() - (sqlc.arg(stale_seconds)::int || ' seconds')::interval
  AND p.disconnected_time IS NULL
ORDER BY p.last_seen_time DESC
LIMIT 1;
