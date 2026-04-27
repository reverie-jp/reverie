package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *repositoryImpl) ListPostLikes(ctx context.Context, postID ulid.ULID, limit int32) ([]*entity.User, error) {
	rows, err := r.q.ListPostLikes(ctx, sqlc.ListPostLikesParams{
		PostID: postID,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}
	users := make([]*entity.User, len(rows))
	for i, u := range rows {
		users[i] = &entity.User{
			ID:          u.ID,
			CustomID:    u.CustomID,
			DisplayName: u.DisplayName,
			IsPrivate:   u.IsPrivate,
			CreateTime:  u.CreateTime,
		}
	}
	return users, nil
}
