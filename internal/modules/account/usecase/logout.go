package usecase

import (
	"context"

	accountrepo "reverie.jp/reverie/internal/modules/account/repository"
)

type Logout struct {
	accountRepo accountrepo.Repository
}

func NewLogout(accountRepo accountrepo.Repository) *Logout {
	return &Logout{accountRepo: accountRepo}
}

func (uc *Logout) Execute(ctx context.Context, input LogoutInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	// Idempotent: missing rows silently succeed, and user_id scoping prevents
	// revoking another user's refresh token even on hash collision / theft.
	return uc.accountRepo.DeleteRefreshTokenByRaw(ctx, input.RefreshToken, input.UserID)
}
