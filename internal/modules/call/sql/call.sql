-- name: CreateCall :exec
INSERT INTO calls (id, host_user_id, visibility)
VALUES ($1, $2, $3);

-- name: ListCallsByIDs :many
SELECT * FROM calls
WHERE id = ANY(sqlc.arg(ids)::text[]);

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
  AND c.end_time IS NULL
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
    display_name = EXCLUDED.display_name,
    -- Fresh rejoin after disconnect: reset host-mute since LiveKit creates
    -- a new unmuted track. Token refresh on an active participant preserves
    -- it (disconnected_time was NULL before upsert).
    muted_by_host = CASE
      WHEN call_participants.disconnected_time IS NOT NULL THEN FALSE
      ELSE call_participants.muted_by_host
    END;

-- name: SetCallParticipantMutedByHost :exec
UPDATE call_participants
SET muted_by_host = TRUE
WHERE call_id = sqlc.arg(call_id)::ulid
  AND participant_identity = sqlc.arg(participant_identity);

-- name: ClearCallParticipantMutedByHost :execrows
UPDATE call_participants
SET muted_by_host = FALSE
WHERE call_id = sqlc.arg(call_id)::ulid
  AND participant_identity = sqlc.arg(participant_identity)
  AND muted_by_host = TRUE;

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

-- name: CreateCallBan :exec
INSERT INTO call_bans (call_id, user_id)
VALUES (sqlc.arg(call_id)::ulid, sqlc.arg(user_id)::ulid)
ON CONFLICT DO NOTHING;

-- name: IsUserBannedFromCall :one
SELECT EXISTS (
  SELECT 1 FROM call_bans
  WHERE call_id = sqlc.arg(call_id)::ulid
    AND user_id = sqlc.arg(user_id)::ulid
) AS banned;

-- name: ListCallBans :many
SELECT * FROM call_bans
WHERE call_id = sqlc.arg(call_id)::ulid
  AND (sqlc.arg(cursor_user_id)::text = '' OR user_id < sqlc.arg(cursor_user_id)::ulid)
ORDER BY user_id DESC
LIMIT sqlc.arg(page_size)::int;

-- name: DeleteCallBan :exec
DELETE FROM call_bans
WHERE call_id = sqlc.arg(call_id)::ulid
  AND user_id = sqlc.arg(user_id)::ulid;

-- name: UpdateCallHost :exec
UPDATE calls
SET host_user_id = sqlc.arg(host_user_id)::ulid, update_time = NOW()
WHERE id = sqlc.arg(id)::ulid;

-- name: MarkAllCallParticipantsDisconnected :exec
UPDATE call_participants
SET disconnected_time = NOW()
WHERE call_id = sqlc.arg(call_id)::ulid
  AND disconnected_time IS NULL;

-- name: MarkCallEnded :exec
UPDATE calls
SET end_time = NOW(), update_time = NOW()
WHERE id = sqlc.arg(id)::ulid
  AND end_time IS NULL;
