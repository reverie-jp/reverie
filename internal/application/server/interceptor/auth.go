package interceptor

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/gen/pb/account/v1/accountv1connect"
	"reverie.jp/reverie/internal/gen/pb/call/v1/callv1connect"
	"reverie.jp/reverie/internal/gen/pb/follow/v1/followv1connect"
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
	followv1connect.FollowServiceListFollowingUsersProcedure:   true,
	followv1connect.FollowServiceListUserFollowersProcedure:    true,
}

type authInterceptor struct {
	jwtManager *jwt.Manager
}

// AuthInterceptor verifies JWT and injects the user ID into context for both
// unary and streaming procedures. Streaming handlers (e.g. EventService) rely
// on this for the subscriber's identity.
func AuthInterceptor(jwtManager *jwt.Manager) connect.Interceptor {
	return &authInterceptor{jwtManager: jwtManager}
}

func (a *authInterceptor) authenticate(ctx context.Context, procedure string, header http.Header) (context.Context, error) {
	if publicProcedures[procedure] {
		return ctx, nil
	}

	rawHeader := header.Get("Authorization")
	optional := optionalAuthProcedures[procedure]
	if rawHeader == "" && optional {
		return ctx, nil
	}

	token, err := extractBearerToken(rawHeader)
	if err != nil {
		return ctx, connect.NewError(connect.CodeUnauthenticated, errors.New("missing or invalid authorization header"))
	}

	claims, err := a.jwtManager.VerifyToken(token)
	if err != nil {
		return ctx, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid or expired token"))
	}

	if claims.TokenType != jwt.TokenTypeAccess {
		return ctx, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token type"))
	}

	userID, err := ulid.Parse(claims.Subject)
	if err != nil {
		return ctx, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token"))
	}

	return ContextWithUserID(ctx, userID), nil
}

func (a *authInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		ctx, err := a.authenticate(ctx, req.Spec().Procedure, req.Header())
		if err != nil {
			return nil, err
		}
		return next(ctx, req)
	}
}

func (a *authInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (a *authInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		slog.Info("streaming request received",
			slog.String("procedure", conn.Spec().Procedure),
			slog.Bool("has_auth_header", conn.RequestHeader().Get("Authorization") != ""),
		)
		ctx, err := a.authenticate(ctx, conn.Spec().Procedure, conn.RequestHeader())
		if err != nil {
			slog.Warn("streaming auth failed",
				slog.String("procedure", conn.Spec().Procedure),
				slog.String("err", err.Error()),
			)
			return err
		}
		return next(ctx, conn)
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
