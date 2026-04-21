package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) GetCall(ctx context.Context, id ulid.ULID) (*entity.Call, error) {
	calls, err := r.ListCallsByIDs(ctx, []ulid.ULID{id})
	if err != nil {
		return nil, err
	}
	if len(calls) == 0 {
		return nil, nil
	}
	return calls[0], nil
}
