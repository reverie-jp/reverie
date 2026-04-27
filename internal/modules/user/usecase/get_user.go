package usecase

import (
	"context"

	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type GetUser struct {
	userGateway usergw.Gateway
}

func NewGetUser(userGateway usergw.Gateway) *GetUser {
	return &GetUser{userGateway: userGateway}
}

func (uc *GetUser) Execute(ctx context.Context, input GetUserInput, requestorID ulid.ULID) (*GetUserOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user, err := resolveUser(ctx, uc.userGateway, input.UserID)
	if err != nil {
		return nil, err
	}

	view, err := uc.userGateway.BuildUserView(ctx, requestorID, user.ID)
	if err != nil {
		return nil, err
	}
	if view == nil {
		return nil, xerrors.ErrUserNotFound
	}

	return &GetUserOutput{View: view}, nil
}
