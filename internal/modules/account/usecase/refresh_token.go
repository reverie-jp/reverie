package usecase

import (
	"context"

	"reverie.jp/reverie/internal/platform/jwt"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type RefreshToken struct {
	jwtManager *jwt.Manager
}

func NewRefreshToken(jwtManager *jwt.Manager) *RefreshToken {
	return &RefreshToken{
		jwtManager: jwtManager,
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

	accessToken, err := uc.jwtManager.GenerateAccessToken(userID)
	if err != nil {
		return nil, xerrors.ErrInvalidRefreshToken.WithCause(err)
	}

	refreshToken, err := uc.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, xerrors.ErrInvalidRefreshToken.WithCause(err)
	}

	return &RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
