package repository

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type ListFollowingTimelineParams struct {
	FollowerID ulid.ULID
	Cursor     *time.Time
	Limit      int32
}

func (r *repositoryImpl) ListFollowingTimeline(ctx context.Context, params ListFollowingTimelineParams) ([]*entity.Post, error) {
	cursor := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	if params.Cursor != nil {
		cursor = *params.Cursor
	}

	limit := params.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	rows, err := r.q.ListFollowingTimeline(ctx, sqlc.ListFollowingTimelineParams{
		FollowerID: params.FollowerID,
		Column2:    cursor,
		Limit:      limit,
	})
	if err != nil {
		return nil, err
	}

	posts := make([]*entity.Post, len(rows))
	for i := range rows {
		posts[i] = mapper.ToPost(&rows[i])
	}
	return posts, nil
}
