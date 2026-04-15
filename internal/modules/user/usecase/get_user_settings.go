package usecase

import (
	"context"

	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type GetUserSettings struct {
	userGateway usergw.Gateway
}

func NewGetUserSettings(userGateway usergw.Gateway) *GetUserSettings {
	return &GetUserSettings{userGateway: userGateway}
}

func (uc *GetUserSettings) Execute(ctx context.Context, userID ulid.ULID) (*GetUserSettingsOutput, error) {
	user, err := uc.userGateway.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, xerrors.ErrUserNotFound
	}

	return &GetUserSettingsOutput{
		IsPrivate: user.IsPrivate,
	}, nil
}
