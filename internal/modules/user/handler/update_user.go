package handler

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"reverie.jp/reverie/internal/application/server/interceptor"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/modules/user/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) UpdateUser(ctx context.Context, req *connect.Request[userv1.UpdateUserRequest]) (*connect.Response[userv1.UpdateUserResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	u := req.Msg.User
	if u == nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("user is required")
	}

	output, err := h.updateUser.Execute(ctx, usecase.UpdateUserInput{
		DisplayName: u.DisplayName,
		Biography:   derefString(u.Biography),
		IsPrivate:   u.IsPrivate,
	}, userID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&userv1.UpdateUserResponse{
		User: &userv1.User{
			Id:          output.ID.String(),
			CustomId:    output.CustomID,
			DisplayName: output.DisplayName,
			Biography:   &output.Biography,
			IsPrivate:   output.IsPrivate,
			IsMe:        true,
			CreateTime:  timestamppb.New(output.CreateTime),
		},
	}), nil
}
