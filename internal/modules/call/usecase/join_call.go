package usecase

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	callrepo "reverie.jp/reverie/internal/modules/call/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/livekit"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type JoinCall struct {
	callRepo    callrepo.Repository
	userGateway usergw.Gateway
	livekit     *livekit.Client
	tokenTTL    time.Duration
}

func NewJoinCall(callRepo callrepo.Repository, userGateway usergw.Gateway, lk *livekit.Client, tokenTTL time.Duration) *JoinCall {
	return &JoinCall{
		callRepo:    callRepo,
		userGateway: userGateway,
		livekit:     lk,
		tokenTTL:    tokenTTL,
	}
}

func (uc *JoinCall) Execute(ctx context.Context, input JoinCallInput) (*JoinCallOutput, error) {
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
	if call.EndTime != nil {
		return nil, xerrors.ErrCallEnded
	}

	if err := uc.checkJoinVisibility(call, input.RequesterID); err != nil {
		return nil, err
	}

	if !input.RequesterID.IsZero() {
		banned, err := uc.callRepo.IsUserBannedFromCall(ctx, call.ID, input.RequesterID)
		if err != nil {
			return nil, xerrors.ErrInternal.WithCause(err)
		}
		if banned {
			return nil, xerrors.ErrCallBanned
		}
	}

	identity, displayName, userID, err := uc.resolveParticipant(ctx, input)
	if err != nil {
		return nil, err
	}

	if err := uc.enforceSingleCall(ctx, input.RequesterID, call.ID); err != nil {
		return nil, err
	}

	if err := uc.callRepo.UpsertCallParticipant(ctx, callrepo.UpsertCallParticipantParams{
		CallID:              call.ID,
		ParticipantIdentity: identity,
		UserID:              userID,
		DisplayName:         displayName,
	}); err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	expireTime := time.Now().Add(uc.tokenTTL)
	token, err := uc.livekit.CreateJoinToken(livekit.JoinTokenParams{
		RoomID:      call.ID.String(),
		Identity:    identity,
		DisplayName: displayName,
		TTL:         uc.tokenTTL,
	})
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	return &JoinCallOutput{
		AccessToken: token,
		URL:         uc.livekit.URL(),
		Identity:    identity,
		ExpireTime:  expireTime,
	}, nil
}

func (uc *JoinCall) checkJoinVisibility(call *entity.Call, requesterID ulid.ULID) error {
	switch call.Visibility {
	case entity.CallVisibilityOpen:
		return nil
	case entity.CallVisibilityUsersOnly:
		if requesterID.IsZero() {
			return xerrors.ErrCallGuestNotAllowed
		}
		return nil
	case entity.CallVisibilityLocked:
		if call.HostUserID == requesterID {
			return nil
		}
		return xerrors.ErrCallLocked
	default:
		return xerrors.ErrCallNotFound
	}
}

func (uc *JoinCall) resolveParticipant(ctx context.Context, input JoinCallInput) (identity, displayName string, userID *ulid.ULID, err error) {
	if input.RequesterID.IsZero() {
		return "guest:" + ulid.New().String(), input.GuestDisplayName, nil, nil
	}

	view, err := uc.userGateway.BuildUserView(ctx, input.RequesterID, input.RequesterID)
	if err != nil {
		return "", "", nil, xerrors.ErrInternal.WithCause(err)
	}
	if view == nil {
		return "", "", nil, xerrors.ErrUserNotFound
	}
	uid := view.User.ID
	return "user:" + uid.String(), view.User.DisplayName, &uid, nil
}

func (uc *JoinCall) enforceSingleCall(ctx context.Context, requesterID, targetCallID ulid.ULID) error {
	if requesterID.IsZero() {
		return nil
	}
	active, err := uc.callRepo.GetActiveCallByUser(ctx, requesterID, entity.ParticipantStaleSeconds)
	if err != nil {
		return xerrors.ErrInternal.WithCause(err)
	}
	if active != nil && active.ID != targetCallID {
		return xerrors.ErrAlreadyInAnotherCall
	}
	return nil
}
