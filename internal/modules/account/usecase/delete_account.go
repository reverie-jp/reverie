package usecase

import (
	"context"

	"reverie.jp/reverie/internal/modules/account/repository"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type DeleteAccount struct {
	repo repository.Repository
}

func NewDeleteAccount(repo repository.Repository) *DeleteAccount {
	return &DeleteAccount{repo: repo}
}

func (uc *DeleteAccount) Execute(ctx context.Context, input DeleteAccountInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	user, err := uc.repo.GetUserByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	if user.CustomID != input.ConfirmCustomID {
		return xerrors.ErrCustomIDMismatch
	}

	return uc.repo.SoftDeleteUser(ctx, input.UserID)
}
