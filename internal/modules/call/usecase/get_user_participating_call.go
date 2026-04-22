package usecase

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	callgw "reverie.jp/reverie/internal/modules/call/gateway"
	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type GetUserParticipatingCall struct {
	callRepo    callrepo.Repository
	userGateway usergw.Gateway
	callGateway callgw.Gateway
}

func NewGetUserParticipatingCall(callRepo callrepo.Repository, userGateway usergw.Gateway, callGateway callgw.Gateway) *GetUserParticipatingCall {
	return &GetUserParticipatingCall{
		callRepo:    callRepo,
		userGateway: userGateway,
		callGateway: callGateway,
	}
}

func (uc *GetUserParticipatingCall) Execute(ctx context.Context, input GetUserParticipatingCallInput) (*GetUserParticipatingCallOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user, err := uc.userGateway.GetUserByCustomID(ctx, input.UserCustomID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if user == nil {
		return nil, xerrors.ErrUserNotFound
	}

	call, err := uc.callRepo.GetActiveCallByUser(ctx, user.ID, entity.ParticipantStaleSeconds)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if call == nil {
		return &GetUserParticipatingCallOutput{}, nil
	}

	participants, err := uc.callRepo.ListCallParticipants(ctx, call.ID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	// Profile viewer has no guest identity context, so pass empty string.
	if err := checkViewVisibility(call, input.RequesterID, "", participants); err != nil {
		return &GetUserParticipatingCallOutput{}, nil
	}

	view, err := uc.callGateway.BuildCallView(ctx, input.RequesterID, call.ID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	return &GetUserParticipatingCallOutput{View: view}, nil
}
