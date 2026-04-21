package usecase

import (
	"context"

	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type HeartbeatCall struct {
	callRepo callrepo.Repository
}

func NewHeartbeatCall(callRepo callrepo.Repository) *HeartbeatCall {
	return &HeartbeatCall{callRepo: callRepo}
}

func (uc *HeartbeatCall) Execute(ctx context.Context, input HeartbeatCallInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	rows, err := uc.callRepo.HeartbeatCallParticipant(ctx, input.CallID, input.Identity)
	if err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	if rows == 0 {
		return xerrors.ErrCallParticipantNotFound
	}
	return nil
}
