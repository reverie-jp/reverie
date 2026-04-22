package usecase

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	callgw "reverie.jp/reverie/internal/modules/call/gateway"
	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type GetCall struct {
	callRepo    callrepo.Repository
	callGateway callgw.Gateway
}

func NewGetCall(callRepo callrepo.Repository, callGateway callgw.Gateway) *GetCall {
	return &GetCall{
		callRepo:    callRepo,
		callGateway: callGateway,
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

	participants, err := uc.callRepo.ListCallParticipants(ctx, call.ID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	if err := checkViewVisibility(call, input.RequesterID, input.ViewerIdentity, participants); err != nil {
		return nil, err
	}

	view, err := uc.callGateway.BuildCallView(ctx, input.RequesterID, call.ID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if view == nil {
		return nil, xerrors.ErrCallNotFound
	}

	participantViews, err := uc.callGateway.BuildListParticipantViews(ctx, input.RequesterID, participants)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	return &GetCallOutput{
		View:         view,
		Participants: participantViews,
	}, nil
}

// checkViewVisibility permits a caller to read call metadata. viewerIdentity
// is the LiveKit identity of the caller if they are currently participating;
// used so that in-progress guests still see LOCKED calls.
func checkViewVisibility(call *entity.Call, requesterID ulid.ULID, viewerIdentity string, participants []*entity.CallParticipant) error {
	switch call.Visibility {
	case entity.CallVisibilityOpen, entity.CallVisibilityUsersOnly:
		return nil
	case entity.CallVisibilityLocked:
		if !requesterID.IsZero() && call.HostUserID == requesterID {
			return nil
		}
		userIdentity := ""
		if !requesterID.IsZero() {
			userIdentity = "user:" + requesterID.String()
		}
		for _, p := range participants {
			if !p.IsCurrentlyConnected() {
				continue
			}
			if userIdentity != "" && p.ParticipantIdentity == userIdentity {
				return nil
			}
			if viewerIdentity != "" && p.ParticipantIdentity == viewerIdentity {
				return nil
			}
		}
		return xerrors.ErrCallNotFound
	default:
		return xerrors.ErrCallNotFound
	}
}
