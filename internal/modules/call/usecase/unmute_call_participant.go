package usecase

import (
	"context"

	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/livekit"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type UnmuteCallParticipant struct {
	callRepo callrepo.Repository
	livekit  *livekit.Client
}

func NewUnmuteCallParticipant(callRepo callrepo.Repository, lk *livekit.Client) *UnmuteCallParticipant {
	return &UnmuteCallParticipant{callRepo: callRepo, livekit: lk}
}

func (uc *UnmuteCallParticipant) Execute(ctx context.Context, input UnmuteCallParticipantInput) error {
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

	// Only allow unmute when muted_by_host is set — self-muted users retain
	// control of their own mic. The conditional UPDATE returns 0 rows if the
	// flag is already false.
	cleared, err := uc.callRepo.ClearCallParticipantMutedByHost(ctx, call.ID, input.Identity)
	if err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	if cleared == 0 {
		return xerrors.ErrCannotUnmuteSelfMuted
	}
	if err := uc.livekit.MuteParticipantMicrophone(ctx, call.ID.String(), input.Identity, false); err != nil {
		// Restore DB so the host can retry.
		_ = uc.callRepo.SetCallParticipantMutedByHost(ctx, call.ID, input.Identity)
		return xerrors.ErrInternal.WithCause(err)
	}
	return nil
}
