package usecase

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type CreateCall struct {
	callRepo    callrepo.Repository
	userGateway usergw.Gateway
}

func NewCreateCall(callRepo callrepo.Repository, userGateway usergw.Gateway) *CreateCall {
	return &CreateCall{
		callRepo:    callRepo,
		userGateway: userGateway,
	}
}

func (uc *CreateCall) Execute(ctx context.Context, input CreateCallInput) (*CreateCallOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	view, err := uc.userGateway.BuildView(ctx, input.RequesterID, input.RequesterID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if view == nil {
		return nil, xerrors.ErrUnauthenticated.WithMessage("session is stale, please log in again")
	}

	callID := ulid.New()
	if err := uc.callRepo.CreateCall(ctx, callrepo.CreateCallParams{
		ID:         callID,
		HostUserID: input.RequesterID,
		Visibility: input.Visibility,
	}); err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	return &CreateCallOutput{
		Call: &entity.Call{
			ID:         callID,
			HostUserID: input.RequesterID,
			Visibility: input.Visibility,
		},
		Host: view,
	}, nil
}
