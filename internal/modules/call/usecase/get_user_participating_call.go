package usecase

import (
	"context"

	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type GetUserParticipatingCall struct {
	callRepo    callrepo.Repository
	userGateway usergw.Gateway
}

func NewGetUserParticipatingCall(callRepo callrepo.Repository, userGateway usergw.Gateway) *GetUserParticipatingCall {
	return &GetUserParticipatingCall{
		callRepo:    callRepo,
		userGateway: userGateway,
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

	call, err := uc.callRepo.GetActiveCallByUser(ctx, user.ID, participantStaleSeconds)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if call == nil {
		return &GetUserParticipatingCallOutput{}, nil
	}

	rows, err := uc.callRepo.ListCallParticipants(ctx, call.ID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	activeSet := buildActiveIdentitySet(rows)
	// Profile viewer has no guest identity context, so pass empty string.
	if err := checkViewVisibility(call, input.RequesterID, "", activeSet); err != nil {
		return &GetUserParticipatingCallOutput{}, nil
	}

	host, err := uc.userGateway.BuildView(ctx, input.RequesterID, call.HostUserID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	return &GetUserParticipatingCallOutput{Call: call, Host: host}, nil
}
