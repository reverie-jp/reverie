package handler

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"

	accountv1 "reverie.jp/reverie/internal/gen/pb/account/v1"
	"reverie.jp/reverie/internal/gen/pb/account/v1/accountv1connect"
	"reverie.jp/reverie/internal/modules/account/usecase"
	"reverie.jp/reverie/internal/platform/ulid"
)

type AccountHandler struct {
	accountv1connect.UnimplementedAccountServiceHandler
	socialLogin  *usecase.SocialLogin
	refreshToken *usecase.RefreshToken
	getAccount   *usecase.GetAccount
	deleteAccount *usecase.DeleteAccount
}

func NewAccountHandler(
	socialLogin *usecase.SocialLogin,
	refreshToken *usecase.RefreshToken,
	getAccount *usecase.GetAccount,
	deleteAccount *usecase.DeleteAccount,
) *AccountHandler {
	return &AccountHandler{
		socialLogin:  socialLogin,
		refreshToken: refreshToken,
		getAccount:   getAccount,
		deleteAccount: deleteAccount,
	}
}

func (h *AccountHandler) SocialLogin(ctx context.Context, req *connect.Request[accountv1.SocialLoginRequest]) (*connect.Response[accountv1.SocialLoginResponse], error) {
	provider, err := toProviderString(req.Msg.Provider)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	output, err := h.socialLogin.Execute(ctx, usecase.SocialLoginInput{
		Provider: provider,
		Code:     req.Msg.Code,
	})
	if err != nil {
		slog.ErrorContext(ctx, "social login failed", slog.String("error", err.Error()))
		return nil, connect.NewError(connect.CodeInternal, errors.New("login failed"))
	}

	return connect.NewResponse(&accountv1.SocialLoginResponse{
		TokenPair: &accountv1.TokenPair{
			AccessToken:  output.AccessToken,
			RefreshToken: output.RefreshToken,
		},
		IsNewAccount: output.IsNewAccount,
	}), nil
}

func (h *AccountHandler) RefreshToken(ctx context.Context, req *connect.Request[accountv1.RefreshTokenRequest]) (*connect.Response[accountv1.RefreshTokenResponse], error) {
	output, err := h.refreshToken.Execute(ctx, usecase.RefreshTokenInput{
		RefreshToken: req.Msg.RefreshToken,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid refresh token"))
	}

	return connect.NewResponse(&accountv1.RefreshTokenResponse{
		TokenPair: &accountv1.TokenPair{
			AccessToken:  output.AccessToken,
			RefreshToken: output.RefreshToken,
		},
	}), nil
}

func (h *AccountHandler) GetAccount(ctx context.Context, req *connect.Request[accountv1.GetAccountRequest]) (*connect.Response[accountv1.GetAccountResponse], error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authentication required"))
	}

	output, err := h.getAccount.Execute(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("account not found"))
		}
		slog.ErrorContext(ctx, "get account failed", slog.String("error", err.Error()))
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to get account"))
	}

	return connect.NewResponse(&accountv1.GetAccountResponse{
		Account: &accountv1.Account{
			Id:          output.ID.String(),
			CustomId:    output.CustomID,
			DisplayName: output.DisplayName,
			AvatarUrl:   output.AvatarURL,
		},
	}), nil
}

func (h *AccountHandler) DeleteAccount(ctx context.Context, req *connect.Request[accountv1.DeleteAccountRequest]) (*connect.Response[accountv1.DeleteAccountResponse], error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authentication required"))
	}

	err = h.deleteAccount.Execute(ctx, usecase.DeleteAccountInput{
		UserID:          userID,
		ConfirmCustomID: req.Msg.ConfirmCustomId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("account not found"))
		}
		slog.ErrorContext(ctx, "delete account failed", slog.String("error", err.Error()))
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to delete account"))
	}

	return connect.NewResponse(&accountv1.DeleteAccountResponse{}), nil
}

func (h *AccountHandler) Logout(ctx context.Context, req *connect.Request[accountv1.LogoutRequest]) (*connect.Response[accountv1.LogoutResponse], error) {
	// With stateless JWTs, logout is handled client-side by discarding tokens.
	// Future: add token blacklist for immediate revocation.
	return connect.NewResponse(&accountv1.LogoutResponse{}), nil
}

func toProviderString(p accountv1.AuthProvider) (string, error) {
	switch p {
	case accountv1.AuthProvider_AUTH_PROVIDER_GOOGLE:
		return "google", nil
	default:
		return "", errors.New("unsupported auth provider")
	}
}

// userIDFromContext extracts the user ID from the context.
// This will be set by an auth interceptor.
func userIDFromContext(ctx context.Context) (ulid.ULID, error) {
	v := ctx.Value(userIDKey{})
	if v == nil {
		return ulid.ULID{}, errors.New("user id not found in context")
	}
	id, ok := v.(ulid.ULID)
	if !ok {
		return ulid.ULID{}, errors.New("invalid user id in context")
	}
	return id, nil
}

type userIDKey struct{}

// ContextWithUserID returns a new context with the user ID set.
func ContextWithUserID(ctx context.Context, userID ulid.ULID) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}
