package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/modules/user/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromGetMyUserRequest(ctx context.Context, _ *connect.Request[userv1.GetMyUserRequest]) (usecase.GetUserInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.GetUserInput{}, xerrors.ErrUnauthenticated
	}

	return usecase.GetUserInput{
		RequesterID: requesterID,
		TargetID:    requesterID,
	}, nil
}

func ToGetMyUserResponse(output *usecase.GetUserOutput) *connect.Response[userv1.GetMyUserResponse] {
	return connect.NewResponse(&userv1.GetMyUserResponse{
		User: ToUser(output.View),
	})
}
