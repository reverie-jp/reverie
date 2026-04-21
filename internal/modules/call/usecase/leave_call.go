package usecase

import (
	"context"

	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/livekit"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type LeaveCall struct {
	callRepo callrepo.Repository
	livekit  *livekit.Client
}

func NewLeaveCall(callRepo callrepo.Repository, lk *livekit.Client) *LeaveCall {
	return &LeaveCall{callRepo: callRepo, livekit: lk}
}

func (uc *LeaveCall) Execute(ctx context.Context, input LeaveCallInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	if _, err := uc.callRepo.MarkCallParticipantDisconnected(ctx, input.CallID, input.Identity); err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	if input.RequesterID.IsZero() {
		return nil
	}
	call, err := uc.callRepo.GetCall(ctx, input.CallID)
	if err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	if call == nil || call.EndTime != nil || call.HostUserID != input.RequesterID {
		return nil
	}
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
