package usecase

import (
	"context"

	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

type GetUser struct {
	userGateway usergw.Gateway
}

func NewGetUser(userGateway usergw.Gateway) *GetUser {
	return &GetUser{userGateway: userGateway}
}

func (uc *GetUser) Execute(ctx context.Context, input GetUserInput, requestorID ulid.ULID) (*GetUserOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user, err := resolveUser(ctx, uc.userGateway, input.UserID)
	if err != nil {
		return nil, err
	}

	isMe := user.ID == requestorID

	followerCount, err := uc.userGateway.CountFollowers(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	followingCount, err := uc.userGateway.CountFollowing(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	var isFollowing, isFollowedBy bool
	if !isMe {
		isFollowing, err = uc.userGateway.IsFollowing(ctx, requestorID, user.ID)
		if err != nil {
			return nil, err
		}
		isFollowedBy, err = uc.userGateway.IsFollowing(ctx, user.ID, requestorID)
		if err != nil {
			return nil, err
		}
	}

	return &GetUserOutput{
		ID:             user.ID,
		CustomID:       user.CustomID,
		DisplayName:    user.DisplayName,
		Biography:      user.Biography,
		IsPrivate:      user.IsPrivate,
		IsMe:           isMe,
		IsFollowing:    isFollowing,
		IsFollowedBy:   isFollowedBy,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
		CreateTime:     user.CreateTime,
	}, nil
}
