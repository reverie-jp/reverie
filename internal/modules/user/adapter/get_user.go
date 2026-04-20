package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/modules/user/usecase"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromGetUserRequest(ctx context.Context, req *connect.Request[userv1.GetUserRequest]) (usecase.GetUserInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.GetUserInput{}, xerrors.ErrUnauthenticated
	}

	targetID, err := ulid.Parse(req.Msg.UserId)
	if err != nil {
		return usecase.GetUserInput{}, xerrors.ErrInvalidArgument.WithMessage("invalid user_id")
	}

	return usecase.GetUserInput{
		RequesterID: requesterID,
		TargetID:    targetID,
	}, nil
}

func ToGetUserResponse(output *usecase.GetUserOutput) *connect.Response[userv1.GetUserResponse] {
	return connect.NewResponse(&userv1.GetUserResponse{
		User: ToUser(output.View),
	})
}
