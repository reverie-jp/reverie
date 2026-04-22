package usecase

import (
	"context"
	"log/slog"

	"reverie.jp/reverie/internal/domain/entity"
	followgw "reverie.jp/reverie/internal/modules/follow/gateway"
	notificationgw "reverie.jp/reverie/internal/modules/notification/gateway"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type UnfollowUser struct {
	followGateway       followgw.Gateway
	userGateway         usergw.Gateway
	notificationGateway notificationgw.Gateway
}

func NewUnfollowUser(followGateway followgw.Gateway, userGateway usergw.Gateway, notificationGateway notificationgw.Gateway) *UnfollowUser {
	return &UnfollowUser{
		followGateway:       followGateway,
		userGateway:         userGateway,
		notificationGateway: notificationGateway,
	}
}

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
	// Clear the "user_followed" notification so re-following later produces
	// a fresh notification ID (otherwise the dedup index returns the stale
	// row and the client de-dupes it out).
	if err := uc.notificationGateway.DeleteByTypeActor(ctx, target.ID, entity.NotificationTypeUserFollowed, input.RequesterID); err != nil {
		slog.Warn("unfollow: failed to delete follow notification",
			slog.String("recipient", target.ID.String()),
			slog.String("actor", input.RequesterID.String()),
			slog.String("err", err.Error()),
		)
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
