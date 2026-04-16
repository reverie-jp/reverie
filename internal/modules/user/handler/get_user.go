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

func (h *Handler) GetUser(ctx context.Context, req *connect.Request[userv1.GetUserRequest]) (*connect.Response[userv1.GetUserResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	output, err := h.getUser.Execute(ctx, usecase.GetUserInput{
		UserID: req.Msg.UserId,
	}, requestorID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&userv1.GetUserResponse{
		User: &userv1.User{
			Id:          output.ID.String(),
			CustomId:    output.CustomID,
			DisplayName: output.DisplayName,
			Biography:   &output.Biography,
			IsPrivate:   output.IsPrivate,
			IsMe:        output.IsMe,
			CreateTime:  timestamppb.New(output.CreateTime),
		},
	}), nil
}
