package usecase

import (
	"context"
	"time"

	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/livekit"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type JoinCall struct {
	livekit     *livekit.Client
	userGateway usergw.Gateway
	tokenTTL    time.Duration
}

func NewJoinCall(lk *livekit.Client, userGateway usergw.Gateway, tokenTTL time.Duration) *JoinCall {
	return &JoinCall{
		livekit:     lk,
		userGateway: userGateway,
		tokenTTL:    tokenTTL,
	}
}

func (uc *JoinCall) Execute(ctx context.Context, input JoinCallInput) (*JoinCallOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	identity, displayName, err := uc.resolveParticipant(ctx, input)
	if err != nil {
		return nil, err
	}

	expireTime := time.Now().Add(uc.tokenTTL)
	token, err := uc.livekit.CreateJoinToken(livekit.JoinTokenParams{
		RoomID:      input.RoomID,
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

func (uc *JoinCall) resolveParticipant(ctx context.Context, input JoinCallInput) (identity, displayName string, err error) {
	if input.RequesterID.IsZero() {
		return "guest:" + ulid.New().String(), input.GuestDisplayName, nil
	}

	view, err := uc.userGateway.BuildView(ctx, input.RequesterID, input.RequesterID)
	if err != nil {
		return "", "", xerrors.ErrInternal.WithCause(err)
	}
	if view == nil {
		return "", "", xerrors.ErrNotFound.WithMessage("user not found")
	}
	return "user:" + input.RequesterID.String(), view.User.DisplayName, nil
}
