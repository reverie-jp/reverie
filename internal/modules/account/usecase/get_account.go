package usecase

import (
	"context"

	"reverie.jp/reverie/internal/modules/account/repository"
)

type GetAccount struct {
	repo repository.Repository
}

func NewGetAccount(repo repository.Repository) *GetAccount {
	return &GetAccount{repo: repo}
}

func (uc *GetAccount) Execute(ctx context.Context, input GetAccountInput) (*GetAccountOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user, err := uc.repo.GetUserByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &GetAccountOutput{
		ID:          user.ID,
		CustomID:    user.CustomID,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		CreateTime:  user.CreateTime,
	}, nil
}
