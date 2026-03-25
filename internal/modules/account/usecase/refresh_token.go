package usecase

import (
	"context"
	"fmt"

	"reverie.jp/reverie/internal/platform/jwt"
	"reverie.jp/reverie/internal/platform/ulid"
)

type RefreshTokenInput struct {
	RefreshToken string
}

type RefreshTokenOutput struct {
	AccessToken  string
	RefreshToken string
}

type RefreshToken struct {
	jwtManager *jwt.Manager
}

func NewRefreshToken(jwtManager *jwt.Manager) *RefreshToken {
	return &RefreshToken{
		jwtManager: jwtManager,
	}
}

func (uc *RefreshToken) Execute(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error) {
	claims, err := uc.jwtManager.VerifyToken(input.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims.TokenType != jwt.TokenTypeRefresh {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	userID, err := ulid.Parse(claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in token: %w", err)
	}

	accessToken, err := uc.jwtManager.GenerateAccessToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := uc.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
