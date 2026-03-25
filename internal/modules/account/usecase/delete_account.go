package usecase

import (
	"context"
	"fmt"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type DeleteAccountInput struct {
	UserID          ulid.ULID
	ConfirmCustomID string
}

type DeleteAccount struct {
	q *sqlc.Queries
}

func NewDeleteAccount(q *sqlc.Queries) *DeleteAccount {
	return &DeleteAccount{q: q}
}

func (uc *DeleteAccount) Execute(ctx context.Context, input DeleteAccountInput) error {
	user, err := uc.q.GetUserByID(ctx, input.UserID.String())
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.CustomID != input.ConfirmCustomID {
		return fmt.Errorf("custom id does not match")
	}

	if err := uc.q.SoftDeleteUser(ctx, input.UserID.String()); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
