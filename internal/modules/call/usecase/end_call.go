package usecase

import (
	"context"

	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/livekit"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type EndCall struct {
	callRepo callrepo.Repository
	livekit  *livekit.Client
}

func NewEndCall(callRepo callrepo.Repository, lk *livekit.Client) *EndCall {
	return &EndCall{callRepo: callRepo, livekit: lk}
}

func (uc *EndCall) Execute(ctx context.Context, input EndCallInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	call, err := uc.callRepo.GetCall(ctx, input.CallID)
	if err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	if call == nil {
		return xerrors.ErrCallNotFound
	}
	if call.HostUserID != input.RequesterID {
		return xerrors.ErrNotCallHost
	}

	// Mark as ended first so that any concurrent JoinCall is rejected even if
	// the LiveKit delete races with a reconnect attempt.
	if err := uc.callRepo.MarkCallEnded(ctx, call.ID); err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	if err := uc.livekit.DeleteRoom(ctx, call.ID.String()); err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	if err := uc.callRepo.MarkAllCallParticipantsDisconnected(ctx, call.ID); err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	return nil
}
