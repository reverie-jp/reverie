package handler

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) GetUserSettings(ctx context.Context, req *connect.Request[userv1.GetUserSettingsRequest]) (*connect.Response[userv1.GetUserSettingsResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	output, err := h.getUserSettings.Execute(ctx, userID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&userv1.GetUserSettingsResponse{
		Settings: &userv1.UserSettings{
			IsPrivate: output.IsPrivate,
		},
	}), nil
}
