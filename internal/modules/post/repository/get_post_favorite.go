package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *repositoryImpl) GetPostFavorite(ctx context.Context, userID ulid.ULID, postID ulid.ULID) (bool, error) {
	_, err := r.q.GetPostFavorite(ctx, sqlc.GetPostFavoriteParams{
		UserID: userID,
		PostID: postID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
