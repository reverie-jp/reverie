package usecase

import (
	"context"

	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type TransferCallHost struct {
	callRepo    callrepo.Repository
	userGateway usergw.Gateway
}

func NewTransferCallHost(callRepo callrepo.Repository, userGateway usergw.Gateway) *TransferCallHost {
	return &TransferCallHost{callRepo: callRepo, userGateway: userGateway}
}

func (uc *TransferCallHost) Execute(ctx context.Context, input TransferCallHostInput) (*TransferCallHostOutput, error) {
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

	newHost, err := uc.userGateway.GetUserByCustomID(ctx, input.NewHostCustomID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if newHost == nil {
		return nil, xerrors.ErrUserNotFound
	}
	if newHost.ID == call.HostUserID {
		return nil, xerrors.ErrInvalidArgument.WithMessage("new host is already the host")
	}

	active, err := uc.callRepo.GetActiveCallByUser(ctx, newHost.ID, participantStaleSeconds)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if active == nil || active.ID != call.ID {
		return nil, xerrors.ErrInvalidArgument.WithMessage("new host must be an active participant of this call")
	}

	if err := uc.callRepo.UpdateCallHost(ctx, call.ID, newHost.ID); err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	call.HostUserID = newHost.ID

	hostView, err := uc.userGateway.BuildView(ctx, input.RequesterID, newHost.ID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	return &TransferCallHostOutput{Call: call, Host: hostView}, nil
}
