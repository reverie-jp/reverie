package usecase

import (
	"context"

	"reverie.jp/reverie/internal/modules/user/gateway"
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

	view, err := uc.userGateway.BuildView(ctx, input.RequesterID, input.TargetID)
	if err != nil {
		return nil, err
	}
	if view == nil {
		return nil, xerrors.ErrUserNotFound
	}

	return &GetUserOutput{View: view}, nil
}
