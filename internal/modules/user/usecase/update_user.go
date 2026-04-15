package usecase

import (
	"context"

	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type UpdateUser struct {
	userGateway usergw.Gateway
}

func NewUpdateUser(userGateway usergw.Gateway) *UpdateUser {
	return &UpdateUser{userGateway: userGateway}
}

func (uc *UpdateUser) Execute(ctx context.Context, input UpdateUserInput, userID ulid.ULID) (*UpdateUserOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user, err := uc.userGateway.UpdateUser(ctx, usergw.UpdateUserParams{
		ID:          userID,
		DisplayName: input.DisplayName,
		Biography:   input.Biography,
		IsPrivate:   input.IsPrivate,
		Birthdate:   input.Birthdate,
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, xerrors.ErrUserNotFound
	}

	return &UpdateUserOutput{
		ID:          user.ID,
		CustomID:    user.CustomID,
		DisplayName: user.DisplayName,
		Biography:   user.Biography,
		IsPrivate:   user.IsPrivate,
		CreateTime:  user.CreateTime,
		UpdateTime:  user.UpdateTime,
	}, nil
}
