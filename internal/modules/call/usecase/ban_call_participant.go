package usecase

import (
	"context"
	"strings"

	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/livekit"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type BanCallParticipant struct {
	callRepo callrepo.Repository
	livekit  *livekit.Client
}

func NewBanCallParticipant(callRepo callrepo.Repository, lk *livekit.Client) *BanCallParticipant {
	return &BanCallParticipant{callRepo: callRepo, livekit: lk}
}

func (uc *BanCallParticipant) Execute(ctx context.Context, input BanCallParticipantInput) error {
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

	userID, ok := parseUserIdentity(input.Identity)
	if !ok {
		return xerrors.ErrCannotBanGuest
	}

	if err := uc.callRepo.CreateCallBan(ctx, call.ID, userID); err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	if err := uc.livekit.RemoveParticipant(ctx, call.ID.String(), input.Identity); err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	if _, err := uc.callRepo.MarkCallParticipantDisconnected(ctx, call.ID, input.Identity); err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	return nil
}

// parseUserIdentity extracts the ULID from a "user:<ULID>" identity. Returns
// ok=false for guest identities.
func parseUserIdentity(identity string) (ulid.ULID, bool) {
	const prefix = "user:"
	if !strings.HasPrefix(identity, prefix) {
		return ulid.ULID{}, false
	}
	id, err := ulid.Parse(identity[len(prefix):])
	if err != nil {
		return ulid.ULID{}, false
	}
	return id, true
}
