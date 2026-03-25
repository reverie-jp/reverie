package handler

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/gen/pb/account/v1/accountv1connect"
	"reverie.jp/reverie/internal/platform/jwt"
	"reverie.jp/reverie/internal/platform/ulid"
)

// publicProcedures lists RPC procedures that do not require authentication.
var publicProcedures = map[string]bool{
	accountv1connect.AccountServiceSocialLoginProcedure:  true,
	accountv1connect.AccountServiceRefreshTokenProcedure: true,
}

func NewAuthInterceptor(jwtManager *jwt.Manager) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if publicProcedures[req.Spec().Procedure] {
				return next(ctx, req)
			}

			token, err := extractBearerToken(req.Header().Get("Authorization"))
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing or invalid authorization header"))
			}

			claims, err := jwtManager.VerifyToken(token)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid or expired token"))
			}

			if claims.TokenType != jwt.TokenTypeAccess {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token type"))
			}

			userID, err := ulid.Parse(claims.Subject)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token"))
			}

			ctx = ContextWithUserID(ctx, userID)
			return next(ctx, req)
		}
	}
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("empty authorization header")
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", errors.New("invalid authorization header format")
	}
	return parts[1], nil
}
