package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *repositoryImpl) DeletePost(ctx context.Context, postID ulid.ULID, authorID ulid.ULID) error {
	return r.q.DeletePost(ctx, sqlc.DeletePostParams{
		ID:       postID,
		AuthorID: authorID,
	})
}
