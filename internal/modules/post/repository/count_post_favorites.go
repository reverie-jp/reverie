package repository

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *repositoryImpl) CountPostFavorites(ctx context.Context, postID ulid.ULID) (int64, error) {
	return r.q.CountPostFavorites(ctx, postID)
}
