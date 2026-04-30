package usecase

import (
	"context"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

type ListPostLikes struct {
	gateway postgw.Gateway
}

func NewListPostLikes(gateway postgw.Gateway) *ListPostLikes {
	return &ListPostLikes{gateway: gateway}
}

func (uc *ListPostLikes) Execute(ctx context.Context, authorCustomID string, shortID string, requestorID ulid.ULID, limit int32) ([]*usergw.UserView, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return uc.gateway.ListPostLikes(ctx, authorCustomID, shortID, requestorID, limit)
}
