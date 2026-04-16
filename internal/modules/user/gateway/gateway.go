package gateway

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/modules/user/repository"
	"reverie.jp/reverie/internal/platform/ulid"
)

type UserView struct {
	User *entity.User
}

type CreateUserParams struct {
	ID          ulid.ULID
	CustomID    string
	DisplayName string
}

type UpdateUserParams struct {
	ID          ulid.ULID
	DisplayName string
	Biography   string
	IsPrivate   bool
	Birthdate   *time.Time
}

type Gateway interface {
	GetUserByID(ctx context.Context, id ulid.ULID) (*entity.User, error)
	GetUserByCustomID(ctx context.Context, customID string) (*entity.User, error)
	ListUsersByIDs(ctx context.Context, ids []ulid.ULID) ([]*entity.User, error)
	CreateUser(ctx context.Context, params CreateUserParams) error
	UpdateUser(ctx context.Context, params UpdateUserParams) (*entity.User, error)
	DeleteUser(ctx context.Context, id ulid.ULID) error
	BuildView(ctx context.Context, id ulid.ULID) (*UserView, error)
	BuildListViews(ctx context.Context, ids []ulid.ULID) ([]*UserView, error)
}

type gatewayImpl struct {
	repo repository.Repository
}

func New(q sqlc.Querier) Gateway {
	return &gatewayImpl{repo: repository.New(q)}
}
