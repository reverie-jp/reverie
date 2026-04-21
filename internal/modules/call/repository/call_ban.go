package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) CreateCallBan(ctx context.Context, callID, userID ulid.ULID) error {
	return r.q.CreateCallBan(ctx, sqlc.CreateCallBanParams{
		CallID: callID,
		UserID: userID,
	})
}

func (r *RepositoryImpl) IsUserBannedFromCall(ctx context.Context, callID, userID ulid.ULID) (bool, error) {
	return r.q.IsUserBannedFromCall(ctx, sqlc.IsUserBannedFromCallParams{
		CallID: callID,
		UserID: userID,
	})
}
