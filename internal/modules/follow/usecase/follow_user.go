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

type FollowUser struct {
	followGateway       followgw.Gateway
	userGateway         usergw.Gateway
	notificationGateway notificationgw.Gateway
}

func NewFollowUser(followGateway followgw.Gateway, userGateway usergw.Gateway, notificationGateway notificationgw.Gateway) *FollowUser {
	return &FollowUser{
		followGateway:       followGateway,
		userGateway:         userGateway,
		notificationGateway: notificationGateway,
	}
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
	// Notification write is best-effort relative to the follow action — a
	// failed insert must not roll back the follow. Log and move on.
	follower := input.RequesterID
	if _, err := uc.notificationGateway.Create(ctx, notificationgw.CreateParams{
		RecipientUserID: target.ID,
		Type:            entity.NotificationTypeUserFollowed,
		ActorUserID:     &follower,
	}); err != nil {
		slog.Warn("follow: failed to create notification",
			slog.String("recipient", target.ID.String()),
			slog.String("actor", follower.String()),
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
	return &FollowUserOutput{View: view}, nil
}
