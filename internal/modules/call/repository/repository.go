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
	Title      string
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
	ListActivePublicCallIDs(ctx context.Context, includeUsersOnly bool, staleSeconds int32, cursorID string, pageSize int32) ([]ulid.ULID, error)
	ListActiveCallIDsForFollower(ctx context.Context, followerID ulid.ULID, staleSeconds int32, cursorID string, pageSize int32) ([]ulid.ULID, error)
	GetActiveCallByUser(ctx context.Context, userID ulid.ULID, staleSeconds int32) (*entity.Call, error)
	UpsertCallParticipant(ctx context.Context, params UpsertCallParticipantParams) error
	HeartbeatCallParticipant(ctx context.Context, callID ulid.ULID, identity string) (int64, error)
	MarkCallParticipantDisconnected(ctx context.Context, callID ulid.ULID, identity string) (int64, error)
	SetCallParticipantMutedByHost(ctx context.Context, callID ulid.ULID, identity string) error
	ClearCallParticipantMutedByHost(ctx context.Context, callID ulid.ULID, identity string) (int64, error)
	ListCallParticipants(ctx context.Context, callID ulid.ULID) ([]*entity.CallParticipant, error)
	// ListActiveParticipantsByCallIDs returns (call_id → []participant) for
	// currently-connected participants (auth + guests) across the given
	// calls. Ordered by first_join_time within each call.
	ListActiveParticipantsByCallIDs(ctx context.Context, callIDs []ulid.ULID, staleSeconds int32) (map[ulid.ULID][]*entity.CallParticipant, error)
	CreateCallBan(ctx context.Context, callID, userID ulid.ULID) error
	IsUserBannedFromCall(ctx context.Context, callID, userID ulid.ULID) (bool, error)
	ListCallBans(ctx context.Context, callID ulid.ULID, cursorUserID string, pageSize int32) ([]*entity.CallBan, error)
	DeleteCallBan(ctx context.Context, callID, userID ulid.ULID) error
	UpdateCallHost(ctx context.Context, callID, hostUserID ulid.ULID) error
	MarkAllCallParticipantsDisconnected(ctx context.Context, callID ulid.ULID) error
	MarkCallEnded(ctx context.Context, callID ulid.ULID) error
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}
