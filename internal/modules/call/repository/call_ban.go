package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
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

func (r *RepositoryImpl) ListCallBans(ctx context.Context, callID ulid.ULID, cursorUserID string, pageSize int32) ([]*entity.CallBan, error) {
	rows, err := r.q.ListCallBans(ctx, sqlc.ListCallBansParams{
		CallID:       callID,
		CursorUserID: cursorUserID,
		PageSize:     pageSize,
	})
	if err != nil {
		return nil, err
	}
	bans := make([]*entity.CallBan, len(rows))
	for i := range rows {
		bans[i] = mapper.ToCallBan(&rows[i])
	}
	return bans, nil
}

func (r *RepositoryImpl) DeleteCallBan(ctx context.Context, callID, userID ulid.ULID) error {
	return r.q.DeleteCallBan(ctx, sqlc.DeleteCallBanParams{
		CallID: callID,
		UserID: userID,
	})
}

func (r *RepositoryImpl) UpdateCallHost(ctx context.Context, callID, hostUserID ulid.ULID) error {
	return r.q.UpdateCallHost(ctx, sqlc.UpdateCallHostParams{
		ID:         callID,
		HostUserID: hostUserID,
	})
}

func (r *RepositoryImpl) MarkAllCallParticipantsDisconnected(ctx context.Context, callID ulid.ULID) error {
	return r.q.MarkAllCallParticipantsDisconnected(ctx, callID)
}

func (r *RepositoryImpl) MarkCallEnded(ctx context.Context, callID ulid.ULID) error {
	return r.q.MarkCallEnded(ctx, callID)
}
