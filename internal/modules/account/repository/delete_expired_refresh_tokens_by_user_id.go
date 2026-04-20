package repository

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) DeleteExpiredRefreshTokensByUserID(ctx context.Context, userID ulid.ULID) error {
	return r.q.DeleteExpiredRefreshTokensByUserID(ctx, userID.String())
}
