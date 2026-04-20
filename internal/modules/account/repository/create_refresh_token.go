package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type CreateRefreshTokenParams struct {
	UserID     ulid.ULID
	RawToken   string
	ExpireTime time.Time
}

func (r *RepositoryImpl) CreateRefreshToken(ctx context.Context, params CreateRefreshTokenParams) error {
	return r.q.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		ID:         ulid.New(),
		UserID:     params.UserID,
		TokenHash:  hashRefreshToken(params.RawToken),
		ExpireTime: pgtype.Timestamptz{Time: params.ExpireTime, Valid: true},
	})
}
