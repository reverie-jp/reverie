package usecase

import (
	"context"

	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type DeleteAccount struct {
	userGateway usergw.Gateway
}

func NewDeleteAccount(userGateway usergw.Gateway) *DeleteAccount {
	return &DeleteAccount{userGateway: userGateway}
}

func (uc *DeleteAccount) Execute(ctx context.Context, input DeleteAccountInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	users, err := uc.userGateway.ListUsersByIDs(ctx, []ulid.ULID{input.UserID})
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return xerrors.ErrAccountNotFound
	}

	if users[0].CustomID != input.ConfirmCustomID {
		return xerrors.ErrCustomIDMismatch
	}

	return uc.userGateway.DeleteUser(ctx, input.UserID)
}
