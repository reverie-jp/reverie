package repository

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *repositoryImpl) CountPostReposts(ctx context.Context, postID ulid.ULID) (int64, error) {
	return r.q.CountPostReposts(ctx, &postID)
}
