package usecase

import (
	"context"
	"encoding/base64"
	"time"

	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

type ListFollowingInput struct {
	UserID    string
	PageToken string
	PageSize  int32
}

type ListFollowing struct {
	userGateway usergw.Gateway
}

func NewListFollowing(userGateway usergw.Gateway) *ListFollowing {
	return &ListFollowing{userGateway: userGateway}
}

func (uc *ListFollowing) Execute(ctx context.Context, input ListFollowingInput, requestorID ulid.ULID) ([]*usergw.UserView, error) {
	target, err := resolveUser(ctx, uc.userGateway, input.UserID)
	if err != nil {
		return nil, err
	}

	var cursor *time.Time
	if input.PageToken != "" {
		if raw, err := base64.StdEncoding.DecodeString(input.PageToken); err == nil {
			if t, err := time.Parse("2006-01-02T15:04:05.999999999Z", string(raw)); err == nil {
				cursor = &t
			}
		}
	}

	limit := input.PageSize
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	users, err := uc.userGateway.ListFollowing(ctx, target.ID, cursor, limit)
	if err != nil {
		return nil, err
	}

	ids := make([]ulid.ULID, len(users))
	for i, u := range users {
		ids[i] = u.ID
	}

	return uc.userGateway.BuildListUserViews(ctx, requestorID, ids)
}
