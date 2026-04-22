package usecase

import (
	"context"
	"log/slog"

	"reverie.jp/reverie/internal/domain/entity"
	callgw "reverie.jp/reverie/internal/modules/call/gateway"
	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	followgw "reverie.jp/reverie/internal/modules/follow/gateway"
	notificationgw "reverie.jp/reverie/internal/modules/notification/gateway"
	"reverie.jp/reverie/internal/platform/resourcename"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type CreateCall struct {
	callRepo            callrepo.Repository
	callGateway         callgw.Gateway
	followGateway       followgw.Gateway
	notificationGateway notificationgw.Gateway
}

func NewCreateCall(
	callRepo callrepo.Repository,
	callGateway callgw.Gateway,
	followGateway followgw.Gateway,
	notificationGateway notificationgw.Gateway,
) *CreateCall {
	return &CreateCall{
		callRepo:            callRepo,
		callGateway:         callGateway,
		followGateway:       followGateway,
		notificationGateway: notificationGateway,
	}
}

func (uc *CreateCall) Execute(ctx context.Context, input CreateCallInput) (*CreateCallOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	callID := ulid.New()
	if err := uc.callRepo.CreateCall(ctx, callrepo.CreateCallParams{
		ID:         callID,
		HostUserID: input.RequesterID,
		Visibility: input.Visibility,
	}); err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	view, err := uc.callGateway.BuildCallView(ctx, input.RequesterID, callID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if view == nil {
		return nil, xerrors.ErrInternal.WithMessage("call disappeared after creation")
	}

	// Fan-out to followers only for discoverable calls. LOCKED is URL-only
	// by design — notifying followers would defeat that. Best-effort: any
	// failure is logged without blocking call creation.
	if input.Visibility != entity.CallVisibilityLocked {
		uc.notifyFollowers(ctx, input.RequesterID, callID)
	}

	return &CreateCallOutput{View: view}, nil
}

func (uc *CreateCall) notifyFollowers(ctx context.Context, hostID, callID ulid.ULID) {
	followerIDs, err := uc.followGateway.ListAllFollowerIDs(ctx, hostID)
	if err != nil {
		slog.Warn("create_call: list followers failed",
			slog.String("host", hostID.String()),
			slog.String("err", err.Error()),
		)
		return
	}
	if len(followerIDs) == 0 {
		return
	}
	actor := hostID
	if err := uc.notificationGateway.FanOutCreate(ctx, notificationgw.FanOutParams{
		RecipientUserIDs: followerIDs,
		Type:             entity.NotificationTypeFollowingUserCallStarted,
		ActorUserID:      &actor,
		ResourceName:     resourcename.FormatCall(callID),
	}); err != nil {
		slog.Warn("create_call: fan-out notifications failed",
			slog.String("host", hostID.String()),
			slog.String("err", err.Error()),
		)
	}
}
