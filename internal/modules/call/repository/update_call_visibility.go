package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) UpdateCallVisibility(ctx context.Context, id ulid.ULID, visibility entity.CallVisibility) error {
	return r.q.UpdateCallVisibility(ctx, sqlc.UpdateCallVisibilityParams{
		ID:         id,
		Visibility: sqlc.CallVisibility(visibility),
	})
}
