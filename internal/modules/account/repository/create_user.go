package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type CreateUserParams struct {
	ID          ulid.ULID
	CustomID    string
	DisplayName string
	AvatarURL   *string
}

func (r *RepositoryImpl) CreateUser(ctx context.Context, params CreateUserParams) error {
	now := time.Now()
	err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		ID:          params.ID,
		CustomID:    params.CustomID,
		DisplayName: params.DisplayName,
		AvatarUrl:   params.AvatarURL,
		CreateTime:  pgtype.Timestamptz{Time: now, Valid: true},
	})
	if err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	return nil
}
