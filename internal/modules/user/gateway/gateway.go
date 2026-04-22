package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	followgw "reverie.jp/reverie/internal/modules/follow/gateway"
	"reverie.jp/reverie/internal/modules/user/repository"
	"reverie.jp/reverie/internal/platform/ulid"
)

type UserView struct {
	User         *entity.User
	IsMe         bool
	IsFollowing  bool
	IsFollowedBy bool
	IsOnline     bool
}

type CreateUserParams struct {
	ID          ulid.ULID
	CustomID    string
	DisplayName string
	AvatarURL   *string
}

type Gateway interface {
	ListUsersByIDs(ctx context.Context, ids []ulid.ULID) ([]*entity.User, error)
	GetUserByCustomID(ctx context.Context, customID string) (*entity.User, error)
	CreateUser(ctx context.Context, params CreateUserParams) error
	DeleteUser(ctx context.Context, id ulid.ULID) error
	// UpdateLastSeen bumps users.last_seen_time to NOW(). Called by the
	// PresenceService.Heartbeat RPC every ~30s while the user is active.
	UpdateLastSeen(ctx context.Context, id ulid.ULID) error
	BuildUserView(ctx context.Context, requesterID, id ulid.ULID) (*UserView, error)
	BuildListUserViews(ctx context.Context, requesterID ulid.ULID, ids []ulid.ULID) ([]*UserView, error)
}

type gatewayImpl struct {
	repo          repository.Repository
	followGateway followgw.Gateway
}

func New(q sqlc.Querier, followGateway followgw.Gateway) Gateway {
	return &gatewayImpl{
		repo:          repository.New(q),
		followGateway: followGateway,
	}
}
