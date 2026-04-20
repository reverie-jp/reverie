package usecase

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/application/transaction"
	"reverie.jp/reverie/internal/gen/sqlc"
	accountrepo "reverie.jp/reverie/internal/modules/account/repository"
	"reverie.jp/reverie/internal/platform/jwt"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type RefreshToken struct {
	accountRepo accountrepo.Repository
	tx          transaction.Runner
	jwtManager  *jwt.Manager
}

func NewRefreshToken(accountRepo accountrepo.Repository, tx transaction.Runner, jwtManager *jwt.Manager) *RefreshToken {
	return &RefreshToken{
		accountRepo: accountRepo,
		tx:          tx,
		jwtManager:  jwtManager,
	}
}

func (uc *RefreshToken) Execute(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	claims, err := uc.jwtManager.VerifyToken(input.RefreshToken)
	if err != nil {
		return nil, xerrors.ErrInvalidRefreshToken.WithCause(err)
	}
	if claims.TokenType != jwt.TokenTypeRefresh {
		return nil, xerrors.ErrInvalidRefreshToken
	}

	userID, err := ulid.Parse(claims.Subject)
	if err != nil {
		return nil, xerrors.ErrInvalidRefreshToken.WithCause(err)
	}

	// Opportunistic cleanup of this user's already-expired refresh tokens.
	if err := uc.accountRepo.DeleteExpiredRefreshTokensByUserID(ctx, userID); err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	record, err := uc.accountRepo.GetRefreshTokenByRaw(ctx, input.RefreshToken)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	if record == nil || record.UserID != userID || time.Now().After(record.ExpireTime) {
		return nil, xerrors.ErrInvalidRefreshToken
	}

	newAccessToken, err := uc.jwtManager.GenerateAccessToken(userID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	newRefreshToken, newExpireTime, err := uc.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	err = uc.tx.WithTx(ctx, func(q sqlc.Querier) error {
		txRepo := accountrepo.New(q)
		if err := txRepo.DeleteRefreshTokenByRaw(ctx, input.RefreshToken, userID); err != nil {
			return err
		}
		return txRepo.CreateRefreshToken(ctx, accountrepo.CreateRefreshTokenParams{
			UserID:     userID,
			RawToken:   newRefreshToken,
			ExpireTime: newExpireTime,
		})
	})
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	return &RefreshTokenOutput{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
