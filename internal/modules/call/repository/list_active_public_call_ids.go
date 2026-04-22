package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) ListActivePublicCallIDs(ctx context.Context, includeUsersOnly bool, staleSeconds int32, cursorID string, pageSize int32) ([]ulid.ULID, error) {
	return r.q.ListActivePublicCallIDs(ctx, sqlc.ListActivePublicCallIDsParams{
		IncludeUsersOnly: includeUsersOnly,
		StaleSeconds:     staleSeconds,
		CursorID:         cursorID,
		PageSize:         pageSize,
	})
}
