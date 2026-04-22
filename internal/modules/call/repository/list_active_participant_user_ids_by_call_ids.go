package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) ListActiveParticipantsByCallIDs(ctx context.Context, callIDs []ulid.ULID, staleSeconds int32) (map[ulid.ULID][]*entity.CallParticipant, error) {
	if len(callIDs) == 0 {
		return map[ulid.ULID][]*entity.CallParticipant{}, nil
	}

	strIDs := make([]string, len(callIDs))
	for i, id := range callIDs {
		strIDs[i] = id.String()
	}

	rows, err := r.q.ListActiveParticipantsByCallIDs(ctx, sqlc.ListActiveParticipantsByCallIDsParams{
		CallIds:      strIDs,
		StaleSeconds: staleSeconds,
	})
	if err != nil {
		return nil, err
	}

	out := make(map[ulid.ULID][]*entity.CallParticipant, len(callIDs))
	for i := range rows {
		p := mapper.ToCallParticipant(&rows[i])
		if p == nil {
			continue
		}
		out[p.CallID] = append(out[p.CallID], p)
	}
	return out, nil
}
