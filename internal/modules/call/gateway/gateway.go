package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/modules/call/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

type CallView struct {
	Call *entity.Call
	Host *usergw.UserView
}

type Gateway interface {
	BuildView(ctx context.Context, requesterID, callID ulid.ULID) (*CallView, error)
	BuildListViews(ctx context.Context, requesterID ulid.ULID, callIDs []ulid.ULID) ([]*CallView, error)
}

type gatewayImpl struct {
	repo        repository.Repository
	userGateway usergw.Gateway
}

func New(repo repository.Repository, userGateway usergw.Gateway) Gateway {
	return &gatewayImpl{
		repo:        repo,
		userGateway: userGateway,
	}
}
