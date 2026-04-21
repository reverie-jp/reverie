package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) ListCallParticipants(ctx context.Context, callID ulid.ULID) ([]*entity.CallParticipant, error) {
	rows, err := r.q.ListCallParticipants(ctx, callID)
	if err != nil {
		return nil, err
	}
	participants := make([]*entity.CallParticipant, len(rows))
	for i := range rows {
		participants[i] = mapper.ToCallParticipant(&rows[i])
	}
	return participants, nil
}
