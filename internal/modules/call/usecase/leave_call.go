package usecase

import (
	"context"

	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type LeaveCall struct {
	callRepo callrepo.Repository
}

func NewLeaveCall(callRepo callrepo.Repository) *LeaveCall {
	return &LeaveCall{callRepo: callRepo}
}

func (uc *LeaveCall) Execute(ctx context.Context, input LeaveCallInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	if _, err := uc.callRepo.MarkCallParticipantDisconnected(ctx, input.CallID, input.Identity); err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	return nil
}
