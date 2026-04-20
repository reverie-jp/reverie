package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) DeleteRefreshTokenByRaw(ctx context.Context, raw string, userID ulid.ULID) error {
	return r.q.DeleteRefreshTokenByHash(ctx, sqlc.DeleteRefreshTokenByHashParams{
		TokenHash: hashRefreshToken(raw),
		UserID:    userID.String(),
	})
}
