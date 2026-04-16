package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *repositoryImpl) DeletePostFavorite(ctx context.Context, userID ulid.ULID, postID ulid.ULID) error {
	return r.q.DeletePostFavorite(ctx, sqlc.DeletePostFavoriteParams{
		UserID: userID,
		PostID: postID,
	})
}
