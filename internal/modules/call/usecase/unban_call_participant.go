package usecase

import (
	"context"

	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type UnbanCallParticipant struct {
	callRepo callrepo.Repository
}

func NewUnbanCallParticipant(callRepo callrepo.Repository) *UnbanCallParticipant {
	return &UnbanCallParticipant{callRepo: callRepo}
}

func (uc *UnbanCallParticipant) Execute(ctx context.Context, input UnbanCallParticipantInput) error {
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

	if err := uc.callRepo.DeleteCallBan(ctx, call.ID, input.UserID); err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	return nil
}
