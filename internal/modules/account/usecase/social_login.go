package usecase

import (
	"context"
	"crypto/rand"
	"math/big"

	"reverie.jp/reverie/internal/application/transaction"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/modules/account/repository"
	"reverie.jp/reverie/internal/platform/google"
	"reverie.jp/reverie/internal/platform/jwt"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type SocialLogin struct {
	repo       repository.Repository
	tx         transaction.Runner
	googleAuth *google.AuthClient
	jwtManager *jwt.Manager
}

func NewSocialLogin(repo repository.Repository, tx transaction.Runner, googleAuth *google.AuthClient, jwtManager *jwt.Manager) *SocialLogin {
	return &SocialLogin{
		repo:       repo,
		tx:         tx,
		googleAuth: googleAuth,
		jwtManager: jwtManager,
	}
}

func (uc *SocialLogin) Execute(ctx context.Context, input SocialLoginInput) (*SocialLoginOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	userInfo, err := uc.googleAuth.Exchange(ctx, input.Code)
	if err != nil {
		return nil, xerrors.ErrSocialLoginFailed.WithCause(err)
	}

	authProvider, err := uc.repo.GetAuthProviderByProvider(ctx, input.Provider, userInfo.Sub)
	if err != nil {
		return nil, xerrors.ErrSocialLoginFailed.WithCause(err)
	}

	var userID ulid.ULID
	isNewAccount := authProvider == nil

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
		return nil, xerrors.ErrSocialLoginFailed.WithCause(err)
	}

	refreshToken, err := uc.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, xerrors.ErrSocialLoginFailed.WithCause(err)
	}

	return &SocialLoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsNewAccount: isNewAccount,
	}, nil
}

func (uc *SocialLogin) createNewUser(ctx context.Context, userInfo *google.UserInfo) (ulid.ULID, error) {
	userID := ulid.New()

	customID, err := generateCustomID()
	if err != nil {
		return ulid.ULID{}, xerrors.ErrSocialLoginFailed.WithCause(err)
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
		txRepo := repository.NewRepository(q)

		if err := txRepo.CreateUser(ctx, repository.CreateUserParams{
			ID:          userID,
			CustomID:    customID,
			DisplayName: displayName,
			AvatarURL:   avatarURL,
		}); err != nil {
			return err
		}

		if err := txRepo.CreateAuthProvider(ctx, repository.CreateAuthProviderParams{
			UserID:         userID,
			Provider:       "google",
			ProviderUserID: userInfo.Sub,
		}); err != nil {
			return err
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
