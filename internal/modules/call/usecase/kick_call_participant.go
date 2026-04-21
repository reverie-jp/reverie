package usecase

import (
	"context"

	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/livekit"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type KickCallParticipant struct {
	callRepo callrepo.Repository
	livekit  *livekit.Client
}

func NewKickCallParticipant(callRepo callrepo.Repository, lk *livekit.Client) *KickCallParticipant {
	return &KickCallParticipant{callRepo: callRepo, livekit: lk}
}

func (uc *KickCallParticipant) Execute(ctx context.Context, input KickCallParticipantInput) error {
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

	if err := uc.livekit.RemoveParticipant(ctx, call.ID.String(), input.Identity); err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	if _, err := uc.callRepo.MarkCallParticipantDisconnected(ctx, call.ID, input.Identity); err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	return nil
}
