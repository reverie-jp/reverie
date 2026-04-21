package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
)

func (r *RepositoryImpl) UpsertCallParticipant(ctx context.Context, params UpsertCallParticipantParams) error {
	return r.q.UpsertCallParticipant(ctx, sqlc.UpsertCallParticipantParams{
		CallID:              params.CallID,
		ParticipantIdentity: params.ParticipantIdentity,
		UserID:              params.UserID,
		DisplayName:         params.DisplayName,
	})
}
