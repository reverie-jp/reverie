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

func (uc *ListFollowing) Execute(ctx context.Context, input ListFollowingInput, requestorID ulid.ULID) ([]*GetUserOutput, error) {
	target, err := resolveUser(ctx, uc.userGateway, input.UserID)
	if err != nil {
		return nil, err
	}
	followerID := target.ID

	var cursor *time.Time
	if input.PageToken != "" {
		raw, err := base64.StdEncoding.DecodeString(input.PageToken)
		if err == nil {
			t, err := time.Parse("2006-01-02T15:04:05.999999999Z", string(raw))
			if err == nil {
				cursor = &t
			}
		}
	}

	limit := input.PageSize
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	users, err := uc.userGateway.ListFollowing(ctx, followerID, cursor, limit)
	if err != nil {
		return nil, err
	}

	outputs := make([]*GetUserOutput, len(users))
	for i, u := range users {
		isMe := u.ID == requestorID
		var isFollowing, isFollowedBy bool
		if !isMe {
			isFollowing, _ = uc.userGateway.IsFollowing(ctx, requestorID, u.ID)
			isFollowedBy, _ = uc.userGateway.IsFollowing(ctx, u.ID, requestorID)
		}
		followerCount, _ := uc.userGateway.CountFollowers(ctx, u.ID)
		followingCount, _ := uc.userGateway.CountFollowing(ctx, u.ID)
		outputs[i] = &GetUserOutput{
			ID:             u.ID,
			CustomID:       u.CustomID,
			DisplayName:    u.DisplayName,
			Biography:      u.Biography,
			IsPrivate:      u.IsPrivate,
			IsMe:           isMe,
			IsFollowing:    isFollowing,
			IsFollowedBy:   isFollowedBy,
			FollowerCount:  followerCount,
			FollowingCount: followingCount,
			CreateTime:     u.CreateTime,
		}
	}
	return outputs, nil
}
