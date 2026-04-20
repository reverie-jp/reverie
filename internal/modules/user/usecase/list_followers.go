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

func (uc *ListFollowers) Execute(ctx context.Context, userID string, pageToken string, pageSize int32, requestorID ulid.ULID) ([]*GetUserOutput, error) {
	target, err := resolveUser(ctx, uc.userGateway, userID)
	if err != nil {
		return nil, err
	}
	followedID := target.ID

	var cursor *time.Time
	if pageToken != "" {
		raw, err := base64.StdEncoding.DecodeString(pageToken)
		if err == nil {
			t, err := time.Parse("2006-01-02T15:04:05.999999999Z", string(raw))
			if err == nil {
				cursor = &t
			}
		}
	}

	limit := pageSize
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	users, err := uc.userGateway.ListFollowers(ctx, followedID, cursor, limit)
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
