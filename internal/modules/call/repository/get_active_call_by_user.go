package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) GetActiveCallByUser(ctx context.Context, userID ulid.ULID, staleSeconds int32) (*entity.Call, error) {
	row, err := r.q.GetActiveCallByUser(ctx, sqlc.GetActiveCallByUserParams{
		UserID:       userID,
		StaleSeconds: staleSeconds,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToCall(&row), nil
}
