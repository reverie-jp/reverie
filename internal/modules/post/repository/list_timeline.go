package repository

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/gen/sqlc"
)

func (r *repositoryImpl) ListTimeline(ctx context.Context, params ListTimelineParams) ([]*entity.Post, error) {
	cursor := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	if params.Cursor != nil {
		cursor = *params.Cursor
	}

	limit := params.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	rows, err := r.q.ListTimeline(ctx, sqlc.ListTimelineParams{
		Column1: cursor,
		Limit:   limit,
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
