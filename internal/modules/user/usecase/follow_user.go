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

	view, err := uc.userGateway.BuildUserView(ctx, requestorID, target.ID)
	if err != nil {
		return nil, err
	}
	return &GetUserOutput{View: view}, nil
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

	view, err := uc.userGateway.BuildUserView(ctx, requestorID, target.ID)
	if err != nil {
		return nil, err
	}
	return &GetUserOutput{View: view}, nil
}
