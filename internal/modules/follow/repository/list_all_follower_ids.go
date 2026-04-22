package repository

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) ListAllFollowerIDs(ctx context.Context, followeeID ulid.ULID) ([]ulid.ULID, error) {
	return r.q.ListAllFollowerIDs(ctx, followeeID)
}
