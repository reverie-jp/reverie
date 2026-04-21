package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/gen/sqlc"
)

func (r *RepositoryImpl) ListActivePublicCalls(ctx context.Context, staleSeconds int32, cursorID string, pageSize int32) ([]*entity.Call, error) {
	rows, err := r.q.ListActivePublicCalls(ctx, sqlc.ListActivePublicCallsParams{
		StaleSeconds: staleSeconds,
		CursorID:     cursorID,
		PageSize:     pageSize,
	})
	if err != nil {
		return nil, err
	}
	calls := make([]*entity.Call, len(rows))
	for i := range rows {
		calls[i] = mapper.ToCall(&rows[i])
	}
	return calls, nil
}
