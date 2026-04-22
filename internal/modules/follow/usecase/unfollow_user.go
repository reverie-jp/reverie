package usecase

import (
	"context"

	followgw "reverie.jp/reverie/internal/modules/follow/gateway"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type UnfollowUser struct {
	followGateway followgw.Gateway
	userGateway   usergw.Gateway
}

func NewUnfollowUser(followGateway followgw.Gateway, userGateway usergw.Gateway) *UnfollowUser {
	return &UnfollowUser{
		followGateway: followGateway,
		userGateway:   userGateway,
	}
}

// Execute removes the follow edge. Intentionally does NOT delete the
// "user_followed" notification that A's earlier follow created — B's bell
// should keep the historical record. Re-follow spam is prevented by the
// notification gateway's per-(recipient, type, actor) cooldown instead.
func (uc *UnfollowUser) Execute(ctx context.Context, input UnfollowUserInput) (*UnfollowUserOutput, error) {
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
	if err := uc.followGateway.DeleteFollow(ctx, input.RequesterID, target.ID); err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	view, err := uc.userGateway.BuildUserView(ctx, input.RequesterID, target.ID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if view == nil {
		return nil, xerrors.ErrUserNotFound
	}
	return &UnfollowUserOutput{View: view}, nil
}
