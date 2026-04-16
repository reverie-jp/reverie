package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *repositoryImpl) CreatePostFavorite(ctx context.Context, userID ulid.ULID, postID ulid.ULID) error {
	return r.q.CreatePostFavorite(ctx, sqlc.CreatePostFavoriteParams{
		UserID: userID,
		PostID: postID,
	})
}
