package repository

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (r *RepositoryImpl) DeleteUser(ctx context.Context, id ulid.ULID) error {
	err := r.q.DeleteUser(ctx, id.String())
	if err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	return nil
}
