package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"reverie.jp/reverie/internal/application/transaction"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/google"
	"reverie.jp/reverie/internal/platform/jwt"
	"reverie.jp/reverie/internal/platform/ulid"
)

type SocialLoginInput struct {
	Provider string
	Code     string
}

type SocialLoginOutput struct {
	AccessToken  string
	RefreshToken string
	IsNewAccount bool
}

type SocialLogin struct {
	q          *sqlc.Queries
	tx         transaction.Runner
	googleAuth *google.AuthClient
	jwtManager *jwt.Manager
}

func NewSocialLogin(q *sqlc.Queries, tx transaction.Runner, googleAuth *google.AuthClient, jwtManager *jwt.Manager) *SocialLogin {
	return &SocialLogin{
		q:          q,
		tx:         tx,
		googleAuth: googleAuth,
		jwtManager: jwtManager,
	}
}

func (uc *SocialLogin) Execute(ctx context.Context, input SocialLoginInput) (*SocialLoginOutput, error) {
	if input.Provider != "google" {
		return nil, fmt.Errorf("unsupported provider: %s", input.Provider)
	}

	userInfo, err := uc.googleAuth.Exchange(ctx, input.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with google: %w", err)
	}

	// Check if this provider account is already linked.
	authProvider, err := uc.q.GetAuthProviderByProvider(ctx, sqlc.GetAuthProviderByProviderParams{
		Provider:       sqlc.AuthProviderGoogle,
		ProviderUserID: userInfo.Sub,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to get auth provider: %w", err)
	}

	var userID ulid.ULID
	isNewAccount := errors.Is(err, pgx.ErrNoRows)

	if isNewAccount {
		userID, err = uc.createNewUser(ctx, userInfo)
		if err != nil {
			return nil, err
		}
	} else {
		userID = authProvider.UserID
	}

	accessToken, err := uc.jwtManager.GenerateAccessToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := uc.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &SocialLoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsNewAccount: isNewAccount,
	}, nil
}

func (uc *SocialLogin) createNewUser(ctx context.Context, userInfo *google.UserInfo) (ulid.ULID, error) {
	userID := ulid.New()
	now := time.Now()

	customID, err := generateCustomID()
	if err != nil {
		return ulid.ULID{}, fmt.Errorf("failed to generate custom id: %w", err)
	}

	displayName := userInfo.Name
	if displayName == "" {
		displayName = "unknown"
	}

	var avatarURL *string
	if userInfo.Picture != "" {
		avatarURL = &userInfo.Picture
	}

	err = uc.tx.WithTx(ctx, func(q sqlc.Querier) error {
		if err := q.CreateUser(ctx, sqlc.CreateUserParams{
			ID:          userID,
			CustomID:    customID,
			DisplayName: displayName,
			AvatarUrl:   avatarURL,
			CreateTime:  pgtype.Timestamptz{Time: now, Valid: true},
		}); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		if err := q.CreateAuthProvider(ctx, sqlc.CreateAuthProviderParams{
			ID:             ulid.New(),
			UserID:         userID,
			Provider:       sqlc.AuthProviderGoogle,
			ProviderUserID: userInfo.Sub,
			CreateTime:     pgtype.Timestamptz{Time: now, Valid: true},
		}); err != nil {
			return fmt.Errorf("failed to create auth provider: %w", err)
		}

		return nil
	})
	if err != nil {
		return ulid.ULID{}, err
	}

	return userID, nil
}

// generateCustomID generates a random 10-character alphanumeric ID.
func generateCustomID() (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	const length = 10

	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[n.Int64()]
	}
	return string(result), nil
}
