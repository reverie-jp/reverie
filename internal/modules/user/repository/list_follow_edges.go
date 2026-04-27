package repository

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) ListFollowingEdges(ctx context.Context, requesterID ulid.ULID, ids []string) ([]ulid.ULID, error) {
	return r.q.ListFollowingEdges(ctx, requesterID, ids)
}

func (r *RepositoryImpl) ListFollowerEdges(ctx context.Context, requesterID ulid.ULID, ids []string) ([]ulid.ULID, error) {
	return r.q.ListFollowerEdges(ctx, requesterID, ids)
}
