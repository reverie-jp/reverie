package interceptor

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/gen/pb/account/v1/accountv1connect"
	"reverie.jp/reverie/internal/gen/pb/call/v1/callv1connect"
	"reverie.jp/reverie/internal/gen/pb/user/v1/userv1connect"
	"reverie.jp/reverie/internal/platform/jwt"
	"reverie.jp/reverie/internal/platform/ulid"
)

type userIDKey struct{}

// ContextWithUserID returns a new context with the user ID set.
func ContextWithUserID(ctx context.Context, userID ulid.ULID) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// UserIDFromContext extracts the user ID from the context.
func UserIDFromContext(ctx context.Context) (ulid.ULID, bool) {
	v := ctx.Value(userIDKey{})
	if v == nil {
		return ulid.ULID{}, false
	}
	id, ok := v.(ulid.ULID)
	return id, ok
}

// publicProcedures lists RPC procedures that skip authentication entirely.
// The Authorization header, if present, is ignored.
var publicProcedures = map[string]bool{
	accountv1connect.AccountServiceSocialLoginProcedure:  true,
	accountv1connect.AccountServiceRefreshTokenProcedure: true,
}

// optionalAuthProcedures lists procedures where authentication is optional:
// callers may omit the Authorization header (and will proceed without a user
// ID in context), but any header that is present must verify successfully.
var optionalAuthProcedures = map[string]bool{
	callv1connect.CallServiceJoinCallProcedure:                 true,
	callv1connect.CallServiceGetCallProcedure:                  true,
	callv1connect.CallServiceListPublicCallsProcedure:          true,
	callv1connect.CallServiceGetUserParticipatingCallProcedure: true,
	callv1connect.CallServiceHeartbeatCallProcedure:            true,
	callv1connect.CallServiceLeaveCallProcedure:                true,
	userv1connect.UserServiceGetUserProcedure:                  true,
}

func AuthInterceptor(jwtManager *jwt.Manager) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if publicProcedures[req.Spec().Procedure] {
				return next(ctx, req)
			}

			header := req.Header().Get("Authorization")
			optional := optionalAuthProcedures[req.Spec().Procedure]
			if header == "" && optional {
				return next(ctx, req)
			}

			token, err := extractBearerToken(header)
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
