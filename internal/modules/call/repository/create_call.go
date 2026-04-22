package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
)

func (r *RepositoryImpl) CreateCall(ctx context.Context, params CreateCallParams) error {
	return r.q.CreateCall(ctx, sqlc.CreateCallParams{
		ID:         params.ID,
		HostUserID: params.HostUserID,
		Visibility: sqlc.CallVisibility(params.Visibility),
		Title:      params.Title,
	})
}
