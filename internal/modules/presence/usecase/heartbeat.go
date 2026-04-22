package usecase

import (
	"context"

	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type Heartbeat struct {
	userGateway usergw.Gateway
}

func NewHeartbeat(userGateway usergw.Gateway) *Heartbeat {
	return &Heartbeat{userGateway: userGateway}
}

func (uc *Heartbeat) Execute(ctx context.Context, input HeartbeatInput) (*HeartbeatOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := uc.userGateway.UpdateLastSeen(ctx, input.RequesterID); err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	return &HeartbeatOutput{}, nil
}
