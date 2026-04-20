-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (id, user_id, token_hash, expire_time)
VALUES ($1, $2, $3, $4);

-- name: GetRefreshTokenByHash :one
SELECT * FROM refresh_tokens
WHERE token_hash = $1;

-- name: DeleteRefreshTokenByHash :exec
DELETE FROM refresh_tokens
WHERE token_hash = $1 AND user_id = $2;

-- name: DeleteExpiredRefreshTokensByUserID :exec
DELETE FROM refresh_tokens
WHERE user_id = $1 AND expire_time < NOW();
