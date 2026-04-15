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

	targetID, err := ulid.Parse(input.UserID)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid user_id")
	}

	user, err := uc.userGateway.GetUserByID(ctx, targetID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, xerrors.ErrUserNotFound
	}

	return &GetUserOutput{
		ID:          user.ID,
		CustomID:    user.CustomID,
		DisplayName: user.DisplayName,
		Biography:   user.Biography,
		IsPrivate:   user.IsPrivate,
		IsMe:        user.ID == requestorID,
		CreateTime:  user.CreateTime,
	}, nil
}
