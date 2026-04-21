package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) HeartbeatCallParticipant(ctx context.Context, callID ulid.ULID, identity string) (int64, error) {
	return r.q.HeartbeatCallParticipant(ctx, sqlc.HeartbeatCallParticipantParams{
		CallID:              callID,
		ParticipantIdentity: identity,
	})
}

func (r *RepositoryImpl) MarkCallParticipantDisconnected(ctx context.Context, callID ulid.ULID, identity string) (int64, error) {
	return r.q.MarkCallParticipantDisconnected(ctx, sqlc.MarkCallParticipantDisconnectedParams{
		CallID:              callID,
		ParticipantIdentity: identity,
	})
}
