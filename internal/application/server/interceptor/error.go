package interceptor

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/go-playground/validator/v10"

	"reverie.jp/reverie/internal/config"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func ErrorInterceptor(env config.Env) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := next(ctx, req)
			if err == nil {
				return resp, nil
			}

			// Already a Connect error (e.g. from auth interceptor).
			var connectErr *connect.Error
			if errors.As(err, &connectErr) {
				return nil, connectErr
			}

			// Domain error with explicit code.
			var appErr *xerrors.Error
			if errors.As(err, &appErr) {
				return nil, connect.NewError(appErr.ConnectCode, appErr)
			}

			// Validation error from go-playground/validator.
			var validationErrs validator.ValidationErrors
			if errors.As(err, &validationErrs) {
				return nil, connect.NewError(connect.CodeInvalidArgument, validationErrs)
			}

			// Unknown error — log and return generic message.
			slog.ErrorContext(ctx, "unhandled error", slog.String("error", err.Error()))
			if env == config.EnvDevelopment {
				return nil, connect.NewError(connect.CodeInternal, err)
			}

			return nil, connect.NewError(connect.CodeInternal, errors.New("internal server error"))
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
