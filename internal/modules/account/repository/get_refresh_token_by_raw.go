package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
)

func (r *RepositoryImpl) GetRefreshTokenByRaw(ctx context.Context, raw string) (*entity.RefreshToken, error) {
	row, err := r.q.GetRefreshTokenByHash(ctx, hashRefreshToken(raw))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToRefreshToken(&row), nil
}
