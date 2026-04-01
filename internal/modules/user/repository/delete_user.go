package repository

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) DeleteUser(ctx context.Context, id ulid.ULID) error {
	return r.q.DeleteUser(ctx, id.String())
}
