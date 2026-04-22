package usecase

import (
	"context"

	followgw "reverie.jp/reverie/internal/modules/follow/gateway"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type FollowUser struct {
	followGateway followgw.Gateway
	userGateway   usergw.Gateway
}

func NewFollowUser(followGateway followgw.Gateway, userGateway usergw.Gateway) *FollowUser {
	return &FollowUser{followGateway: followGateway, userGateway: userGateway}
}

func (uc *FollowUser) Execute(ctx context.Context, input FollowUserInput) (*FollowUserOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	target, err := uc.userGateway.GetUserByCustomID(ctx, input.TargetCustomID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if target == nil {
		return nil, xerrors.ErrUserNotFound
	}
	if target.ID == input.RequesterID {
		return nil, xerrors.ErrCannotFollowSelf
	}
	if err := uc.followGateway.CreateFollow(ctx, input.RequesterID, target.ID); err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	view, err := uc.userGateway.BuildUserView(ctx, input.RequesterID, target.ID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if view == nil {
		return nil, xerrors.ErrUserNotFound
	}
	return &FollowUserOutput{View: view}, nil
}
