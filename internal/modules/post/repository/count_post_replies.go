package repository

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *repositoryImpl) CountPostReplies(ctx context.Context, postID ulid.ULID) (int64, error) {
	return r.q.CountPostReplies(ctx, &postID)
}
