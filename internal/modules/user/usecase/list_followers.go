package usecase

import (
	"context"
	"encoding/base64"
	"time"

	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

type ListFollowers struct {
	userGateway usergw.Gateway
}

func NewListFollowers(userGateway usergw.Gateway) *ListFollowers {
	return &ListFollowers{userGateway: userGateway}
}

func (uc *ListFollowers) Execute(ctx context.Context, userID string, pageToken string, pageSize int32, requestorID ulid.ULID) ([]*usergw.UserView, error) {
	target, err := resolveUser(ctx, uc.userGateway, userID)
	if err != nil {
		return nil, err
	}

	var cursor *time.Time
	if pageToken != "" {
		if raw, err := base64.StdEncoding.DecodeString(pageToken); err == nil {
			if t, err := time.Parse("2006-01-02T15:04:05.999999999Z", string(raw)); err == nil {
				cursor = &t
			}
		}
	}

	limit := pageSize
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	users, err := uc.userGateway.ListFollowers(ctx, target.ID, cursor, limit)
	if err != nil {
		return nil, err
	}

	ids := make([]ulid.ULID, len(users))
	for i, u := range users {
		ids[i] = u.ID
	}

	return uc.userGateway.BuildListUserViews(ctx, requestorID, ids)
}
