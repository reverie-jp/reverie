package usecase

import (
	"context"

	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type UpdateCall struct {
	callRepo    callrepo.Repository
	userGateway usergw.Gateway
}

func NewUpdateCall(callRepo callrepo.Repository, userGateway usergw.Gateway) *UpdateCall {
	return &UpdateCall{
		callRepo:    callRepo,
		userGateway: userGateway,
	}
}

func (uc *UpdateCall) Execute(ctx context.Context, input UpdateCallInput) (*UpdateCallOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	call, err := uc.callRepo.GetCall(ctx, input.CallID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if call == nil {
		return nil, xerrors.ErrCallNotFound
	}
	if call.HostUserID != input.RequesterID {
		return nil, xerrors.ErrNotCallHost
	}

	if err := uc.callRepo.UpdateCallVisibility(ctx, input.CallID, input.Visibility); err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	host, err := uc.userGateway.BuildView(ctx, input.RequesterID, call.HostUserID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	call.Visibility = input.Visibility
	return &UpdateCallOutput{Call: call, Host: host}, nil
}
