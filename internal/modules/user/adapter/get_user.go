package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/modules/user/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromGetUserRequest(ctx context.Context, req *connect.Request[userv1.GetUserRequest]) (usecase.GetUserInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.GetUserInput{}, xerrors.ErrUnauthenticated
	}

	return usecase.GetUserInput{
		RequesterID:    requesterID,
		TargetCustomID: req.Msg.CustomId,
	}, nil
}

func ToGetUserResponse(output *usecase.GetUserOutput) *connect.Response[userv1.GetUserResponse] {
	return connect.NewResponse(&userv1.GetUserResponse{
		User: ToUser(output.View),
	})
}
