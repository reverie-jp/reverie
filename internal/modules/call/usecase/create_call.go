package usecase

import (
	"context"

	callgw "reverie.jp/reverie/internal/modules/call/gateway"
	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type CreateCall struct {
	callRepo    callrepo.Repository
	callGateway callgw.Gateway
}

func NewCreateCall(callRepo callrepo.Repository, callGateway callgw.Gateway) *CreateCall {
	return &CreateCall{
		callRepo:    callRepo,
		callGateway: callGateway,
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

	return &CreateCallOutput{View: view}, nil
}
