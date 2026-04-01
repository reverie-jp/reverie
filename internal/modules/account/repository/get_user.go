package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (r *RepositoryImpl) GetUserByID(ctx context.Context, id ulid.ULID) (*entity.User, error) {
	row, err := r.q.GetUserByID(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, xerrors.ErrAccountNotFound
		}
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	return mapper.ToUser(&row), nil
}
