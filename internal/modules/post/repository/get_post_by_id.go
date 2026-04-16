package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *repositoryImpl) GetPostByID(ctx context.Context, id ulid.ULID) (*entity.Post, error) {
	row, err := r.q.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToPost(&row), nil
}
