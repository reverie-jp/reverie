package usecase

import (
	"context"

	callgw "reverie.jp/reverie/internal/modules/call/gateway"
	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type UpdateCall struct {
	callRepo    callrepo.Repository
	callGateway callgw.Gateway
}

func NewUpdateCall(callRepo callrepo.Repository, callGateway callgw.Gateway) *UpdateCall {
	return &UpdateCall{
		callRepo:    callRepo,
		callGateway: callGateway,
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

	view, err := uc.callGateway.BuildCallView(ctx, input.RequesterID, input.CallID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if view == nil {
		return nil, xerrors.ErrCallNotFound
	}

	return &UpdateCallOutput{View: view}, nil
}
