package usecase

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

// resolveUser looks up a user by ULID string first, falling back to custom_id.
func resolveUser(ctx context.Context, gw usergw.Gateway, userID string) (*entity.User, error) {
	if id, err := ulid.Parse(userID); err == nil {
		user, err := gw.GetUserByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, xerrors.ErrUserNotFound
		}
		return user, nil
	}
	user, err := gw.GetUserByCustomID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, xerrors.ErrUserNotFound
	}
	return user, nil
}

type FollowUser struct {
	userGateway usergw.Gateway
}

func NewFollowUser(userGateway usergw.Gateway) *FollowUser {
	return &FollowUser{userGateway: userGateway}
}

func (uc *FollowUser) Execute(ctx context.Context, targetUserID string, requestorID ulid.ULID) (*GetUserOutput, error) {
	target, err := resolveUser(ctx, uc.userGateway, targetUserID)
	if err != nil {
		return nil, err
	}
	if target.ID == requestorID {
		return nil, xerrors.ErrInvalidArgument.WithMessage("cannot follow yourself")
	}

	if err := uc.userGateway.FollowUser(ctx, requestorID, target.ID); err != nil {
		return nil, err
	}

	return buildGetUserOutput(ctx, uc.userGateway, target.ID, requestorID)
}

type UnfollowUser struct {
	userGateway usergw.Gateway
}

func NewUnfollowUser(userGateway usergw.Gateway) *UnfollowUser {
	return &UnfollowUser{userGateway: userGateway}
}

func (uc *UnfollowUser) Execute(ctx context.Context, targetUserID string, requestorID ulid.ULID) (*GetUserOutput, error) {
	target, err := resolveUser(ctx, uc.userGateway, targetUserID)
	if err != nil {
		return nil, err
	}

	if err := uc.userGateway.UnfollowUser(ctx, requestorID, target.ID); err != nil {
		return nil, err
	}

	return buildGetUserOutput(ctx, uc.userGateway, target.ID, requestorID)
}

func buildGetUserOutput(ctx context.Context, gw usergw.Gateway, targetID, requestorID ulid.ULID) (*GetUserOutput, error) {
	user, err := gw.GetUserByID(ctx, targetID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, xerrors.ErrUserNotFound  //nolint
	}

	followerCount, err := gw.CountFollowers(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	followingCount, err := gw.CountFollowing(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	isFollowing, err := gw.IsFollowing(ctx, requestorID, user.ID)
	if err != nil {
		return nil, err
	}
	isFollowedBy, err := gw.IsFollowing(ctx, user.ID, requestorID)
	if err != nil {
		return nil, err
	}

	return &GetUserOutput{
		ID:             user.ID,
		CustomID:       user.CustomID,
		DisplayName:    user.DisplayName,
		Biography:      user.Biography,
		IsPrivate:      user.IsPrivate,
		IsMe:           false,
		IsFollowing:    isFollowing,
		IsFollowedBy:   isFollowedBy,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
		CreateTime:     user.CreateTime,
	}, nil
}
