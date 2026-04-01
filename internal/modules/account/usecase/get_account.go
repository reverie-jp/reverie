package usecase

import (
	"context"

	userrepo "reverie.jp/reverie/internal/modules/user/repository"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type GetAccount struct {
	userRepo userrepo.Repository
}

func NewGetAccount(userRepo userrepo.Repository) *GetAccount {
	return &GetAccount{userRepo: userRepo}
}

func (uc *GetAccount) Execute(ctx context.Context, input GetAccountInput) (*GetAccountOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	users, err := uc.userRepo.ListUsersByIDs(ctx, []ulid.ULID{input.UserID})
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
		AvatarURL:   user.AvatarURL,
		CreateTime:  user.CreateTime,
	}, nil
}
