package usecase

import (
	"context"

	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type GetAccount struct {
	userGateway usergw.Gateway
}

func NewGetAccount(userGateway usergw.Gateway) *GetAccount {
	return &GetAccount{userGateway: userGateway}
}

func (uc *GetAccount) Execute(ctx context.Context, input GetAccountInput) (*GetAccountOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	users, err := uc.userGateway.ListUsersByIDs(ctx, []ulid.ULID{input.UserID})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, xerrors.ErrAccountNotFound
	}

	user := users[0]

	return &GetAccountOutput{
		ID:          user.ID,
		CustomID:    user.CustomID,
		DisplayName: user.DisplayName,
		CreateTime:  user.CreateTime,
	}, nil
}
