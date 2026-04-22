package repository

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) UpdateUserLastSeen(ctx context.Context, id ulid.ULID) error {
	return r.q.UpdateUserLastSeen(ctx, id)
}
