package usecase

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type GetCall struct {
	callRepo    callrepo.Repository
	userGateway usergw.Gateway
}

func NewGetCall(callRepo callrepo.Repository, userGateway usergw.Gateway) *GetCall {
	return &GetCall{
		callRepo:    callRepo,
		userGateway: userGateway,
	}
}

func (uc *GetCall) Execute(ctx context.Context, input GetCallInput) (*GetCallOutput, error) {
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

	rows, err := uc.callRepo.ListCallParticipants(ctx, call.ID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	activeSet := buildActiveIdentitySet(rows)
	if err := checkViewVisibility(call, input.RequesterID, input.ViewerIdentity, activeSet); err != nil {
		return nil, err
	}

	userIDs := make([]ulid.ULID, 0, len(rows)+1)
	userIDs = append(userIDs, call.HostUserID)
	for _, p := range rows {
		if p.UserID != nil {
			userIDs = append(userIDs, *p.UserID)
		}
	}
	views, err := uc.userGateway.BuildListViews(ctx, input.RequesterID, userIDs)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	viewByID := make(map[ulid.ULID]*usergw.UserView, len(views))
	for _, v := range views {
		if v != nil && v.User != nil {
			viewByID[v.User.ID] = v
		}
	}

	participants := make([]*CallParticipantView, len(rows))
	for i, p := range rows {
		var view *usergw.UserView
		if p.UserID != nil {
			view = viewByID[*p.UserID]
		}
		participants[i] = &CallParticipantView{
			Participant:          p,
			UserView:             view,
			IsCurrentlyConnected: activeSet[p.ParticipantIdentity],
		}
	}

	return &GetCallOutput{
		Call:         call,
		Host:         viewByID[call.HostUserID],
		Participants: participants,
	}, nil
}

func buildActiveIdentitySet(rows []*entity.CallParticipant) map[string]bool {
	cutoff := time.Now().Add(-participantStaleSeconds * time.Second)
	set := make(map[string]bool, len(rows))
	for _, p := range rows {
		if p.DisconnectedTime == nil && p.LastSeenTime.After(cutoff) {
			set[p.ParticipantIdentity] = true
		}
	}
	return set
}

// checkViewVisibility permits a caller to read call metadata.
// viewerIdentity is the LiveKit identity of the caller if they are currently
// participating; used so that in-progress guests still see LOCKED calls.
func checkViewVisibility(call *entity.Call, requesterID ulid.ULID, viewerIdentity string, activeIdentities map[string]bool) error {
	switch call.Visibility {
	case entity.CallVisibilityOpen:
		return nil
	case entity.CallVisibilityUsersOnly:
		if !requesterID.IsZero() {
			return nil
		}
		// Grandfather guests who were already participating when the host
		// tightened visibility — they stay in, but no new guests can view.
		if viewerIdentity != "" && activeIdentities[viewerIdentity] {
			return nil
		}
		return xerrors.ErrUnauthenticated
	case entity.CallVisibilityLocked:
		if !requesterID.IsZero() && call.HostUserID == requesterID {
			return nil
		}
		if !requesterID.IsZero() && activeIdentities["user:"+requesterID.String()] {
			return nil
		}
		if viewerIdentity != "" && activeIdentities[viewerIdentity] {
			return nil
		}
		return xerrors.ErrCallNotFound
	default:
		return xerrors.ErrCallNotFound
	}
}
