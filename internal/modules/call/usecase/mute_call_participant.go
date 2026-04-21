package usecase

import (
	"context"

	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/livekit"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type MuteCallParticipant struct {
	callRepo callrepo.Repository
	livekit  *livekit.Client
}

func NewMuteCallParticipant(callRepo callrepo.Repository, lk *livekit.Client) *MuteCallParticipant {
	return &MuteCallParticipant{callRepo: callRepo, livekit: lk}
}

func (uc *MuteCallParticipant) Execute(ctx context.Context, input MuteCallParticipantInput) error {
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
	if input.Identity == "user:"+call.HostUserID.String() {
		return xerrors.ErrCannotTargetHost
	}

	if err := uc.livekit.MuteParticipantMicrophone(ctx, call.ID.String(), input.Identity, true); err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	if err := uc.callRepo.SetCallParticipantMutedByHost(ctx, call.ID, input.Identity); err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	return nil
}
