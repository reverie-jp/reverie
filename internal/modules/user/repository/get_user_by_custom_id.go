package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
)

func (r *RepositoryImpl) GetUserByCustomID(ctx context.Context, customID string) (*entity.User, error) {
	row, err := r.q.GetUserByCustomID(ctx, customID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToUser(&row), nil
}
