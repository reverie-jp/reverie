package usecase

import (
	"context"

	userrepo "reverie.jp/reverie/internal/modules/user/repository"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type DeleteAccount struct {
	userRepo userrepo.Repository
}

func NewDeleteAccount(userRepo userrepo.Repository) *DeleteAccount {
	return &DeleteAccount{userRepo: userRepo}
}

func (uc *DeleteAccount) Execute(ctx context.Context, input DeleteAccountInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	users, err := uc.userRepo.ListUsersByIDs(ctx, []ulid.ULID{input.UserID})
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return xerrors.ErrAccountNotFound
	}

	if users[0].CustomID != input.ConfirmCustomID {
		return xerrors.ErrCustomIDMismatch
	}

	return uc.userRepo.DeleteUser(ctx, input.UserID)
}
