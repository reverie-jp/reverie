package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) SetCallParticipantMutedByHost(ctx context.Context, callID ulid.ULID, identity string) error {
	return r.q.SetCallParticipantMutedByHost(ctx, sqlc.SetCallParticipantMutedByHostParams{
		CallID:              callID,
		ParticipantIdentity: identity,
	})
}

func (r *RepositoryImpl) ClearCallParticipantMutedByHost(ctx context.Context, callID ulid.ULID, identity string) (int64, error) {
	return r.q.ClearCallParticipantMutedByHost(ctx, sqlc.ClearCallParticipantMutedByHostParams{
		CallID:              callID,
		ParticipantIdentity: identity,
	})
}
