package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type CreateCallParams struct {
	ID         ulid.ULID
	HostUserID ulid.ULID
	Visibility entity.CallVisibility
}

type UpsertCallParticipantParams struct {
	CallID              ulid.ULID
	ParticipantIdentity string
	UserID              *ulid.ULID
	DisplayName         string
}

type Repository interface {
	CreateCall(ctx context.Context, params CreateCallParams) error
	ListCallsByIDs(ctx context.Context, ids []ulid.ULID) ([]*entity.Call, error)
	GetCall(ctx context.Context, id ulid.ULID) (*entity.Call, error)
	UpdateCallVisibility(ctx context.Context, id ulid.ULID, visibility entity.CallVisibility) error
	ListActivePublicCalls(ctx context.Context, staleSeconds int32, cursorID string, pageSize int32) ([]*entity.Call, error)
	GetActiveCallByUser(ctx context.Context, userID ulid.ULID, staleSeconds int32) (*entity.Call, error)
	UpsertCallParticipant(ctx context.Context, params UpsertCallParticipantParams) error
	HeartbeatCallParticipant(ctx context.Context, callID ulid.ULID, identity string) (int64, error)
	MarkCallParticipantDisconnected(ctx context.Context, callID ulid.ULID, identity string) (int64, error)
	ListCallParticipants(ctx context.Context, callID ulid.ULID) ([]*entity.CallParticipant, error)
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}
