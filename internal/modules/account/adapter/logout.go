package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	accountv1 "reverie.jp/reverie/internal/gen/pb/account/v1"
	"reverie.jp/reverie/internal/modules/account/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromLogoutRequest(ctx context.Context, req *connect.Request[accountv1.LogoutRequest]) (usecase.LogoutInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.LogoutInput{}, xerrors.ErrUnauthenticated
	}
	return usecase.LogoutInput{
		UserID:       userID,
		RefreshToken: req.Msg.RefreshToken,
	}, nil
}

func ToLogoutResponse() *connect.Response[accountv1.LogoutResponse] {
	return connect.NewResponse(&accountv1.LogoutResponse{})
}
