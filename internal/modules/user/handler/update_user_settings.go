package handler

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/modules/user/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) UpdateUserSettings(ctx context.Context, req *connect.Request[userv1.UpdateUserSettingsRequest]) (*connect.Response[userv1.UpdateUserSettingsResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	s := req.Msg.Settings
	if s == nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("settings is required")
	}

	output, err := h.updateUserSettings.Execute(ctx, usecase.UpdateUserSettingsInput{
		IsPrivate: s.IsPrivate,
	}, userID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&userv1.UpdateUserSettingsResponse{
		Settings: &userv1.UserSettings{
			IsPrivate: output.IsPrivate,
		},
	}), nil
}
