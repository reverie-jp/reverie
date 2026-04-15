package usecase

import (
	"context"

	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type UpdateUserSettings struct {
	userGateway usergw.Gateway
}

func NewUpdateUserSettings(userGateway usergw.Gateway) *UpdateUserSettings {
	return &UpdateUserSettings{userGateway: userGateway}
}

func (uc *UpdateUserSettings) Execute(ctx context.Context, input UpdateUserSettingsInput, userID ulid.ULID) (*UpdateUserSettingsOutput, error) {
	user, err := uc.userGateway.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, xerrors.ErrUserNotFound
	}

	updated, err := uc.userGateway.UpdateUser(ctx, usergw.UpdateUserParams{
		ID:          userID,
		DisplayName: user.DisplayName,
		Biography:   user.Biography,
		IsPrivate:   input.IsPrivate,
		Birthdate:   user.Birthdate,
	})
	if err != nil {
		return nil, err
	}

	return &UpdateUserSettingsOutput{
		IsPrivate: updated.IsPrivate,
	}, nil
}
