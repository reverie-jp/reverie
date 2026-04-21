package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) ListCallsByIDs(ctx context.Context, ids []ulid.ULID) ([]*entity.Call, error) {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = id.String()
	}

	rows, err := r.q.ListCallsByIDs(ctx, strIDs)
	if err != nil {
		return nil, err
	}

	calls := make([]*entity.Call, len(rows))
	for i := range rows {
		calls[i] = mapper.ToCall(&rows[i])
	}

	return calls, nil
}
