package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) ListUsersByIDs(ctx context.Context, ids []ulid.ULID) ([]*entity.User, error) {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = id.String()
	}

	rows, err := r.q.ListUsersByIDs(ctx, strIDs)
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, len(rows))
	for i := range rows {
		users[i] = mapper.ToUser(&rows[i])
	}

	return users, nil
}
