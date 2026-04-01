package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type CreateAuthProviderParams struct {
	UserID         ulid.ULID
	Provider       string
	ProviderUserID string
}

func (r *RepositoryImpl) CreateAuthProvider(ctx context.Context, params CreateAuthProviderParams) error {
	now := time.Now()
	err := r.q.CreateAuthProvider(ctx, sqlc.CreateAuthProviderParams{
		ID:             ulid.New(),
		UserID:         params.UserID,
		Provider:       sqlc.AuthProvider(params.Provider),
		ProviderUserID: params.ProviderUserID,
		CreateTime:     pgtype.Timestamptz{Time: now, Valid: true},
	})
	if err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	return nil
}
