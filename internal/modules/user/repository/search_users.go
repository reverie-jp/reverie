package repository

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/gen/sqlc"
)

func (r *RepositoryImpl) SearchUsers(ctx context.Context, query string, cursor *time.Time, limit int32) ([]*entity.User, error) {
	cur := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	if cursor != nil {
		cur = *cursor
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	rows, err := r.q.SearchUsers(ctx, sqlc.SearchUsersParams{
		Query:   query,
		Column2: cur,
		Limit:   limit,
	})
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, len(rows))
	for i := range rows {
		users[i] = mapper.ToUser(&rows[i])
	}
	return users, nil
}
