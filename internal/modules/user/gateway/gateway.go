package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/modules/user/repository"
	"reverie.jp/reverie/internal/platform/ulid"
)

type UserView struct {
	User *entity.User
	IsMe bool
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
	BuildView(ctx context.Context, requesterID, id ulid.ULID) (*UserView, error)
	BuildListViews(ctx context.Context, requesterID ulid.ULID, ids []ulid.ULID) ([]*UserView, error)
}

type gatewayImpl struct {
	repo repository.Repository
}

func New(q sqlc.Querier) Gateway {
	return &gatewayImpl{repo: repository.New(q)}
}
