package usecase

import (
	"context"

	"reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type GetUser struct {
	userGateway gateway.Gateway
}

func NewGetUser(userGateway gateway.Gateway) *GetUser {
	return &GetUser{userGateway: userGateway}
}

func (uc *GetUser) Execute(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	targetID, err := uc.resolveTargetID(ctx, input)
	if err != nil {
		return nil, err
	}

	view, err := uc.userGateway.BuildView(ctx, input.RequesterID, targetID)
	if err != nil {
		return nil, err
	}
	if view == nil {
		return nil, xerrors.ErrUserNotFound
	}

	return &GetUserOutput{View: view}, nil
}

func (uc *GetUser) resolveTargetID(ctx context.Context, input GetUserInput) (ulid.ULID, error) {
	if !input.TargetID.IsZero() {
		return input.TargetID, nil
	}
	user, err := uc.userGateway.GetUserByCustomID(ctx, input.TargetCustomID)
	if err != nil {
		return ulid.ULID{}, err
	}
	if user == nil {
		return ulid.ULID{}, xerrors.ErrUserNotFound
	}
	return user.ID, nil
}
