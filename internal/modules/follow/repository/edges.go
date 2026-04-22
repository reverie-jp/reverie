package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) ListFollowingEdges(ctx context.Context, followerID ulid.ULID, targetIDs []ulid.ULID) ([]ulid.ULID, error) {
	if len(targetIDs) == 0 {
		return []ulid.ULID{}, nil
	}
	strIDs := make([]string, len(targetIDs))
	for i, id := range targetIDs {
		strIDs[i] = id.String()
	}
	return r.q.ListFollowingEdgesForRequester(ctx, sqlc.ListFollowingEdgesForRequesterParams{
		FollowerID:  followerID,
		FolloweeIds: strIDs,
	})
}

func (r *RepositoryImpl) ListFollowerEdges(ctx context.Context, followeeID ulid.ULID, targetIDs []ulid.ULID) ([]ulid.ULID, error) {
	if len(targetIDs) == 0 {
		return []ulid.ULID{}, nil
	}
	strIDs := make([]string, len(targetIDs))
	for i, id := range targetIDs {
		strIDs[i] = id.String()
	}
	return r.q.ListFollowerEdgesForRequester(ctx, sqlc.ListFollowerEdgesForRequesterParams{
		FolloweeID:  followeeID,
		FollowerIds: strIDs,
	})
}
