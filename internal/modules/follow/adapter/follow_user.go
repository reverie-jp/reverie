package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	followv1 "reverie.jp/reverie/internal/gen/pb/follow/v1"
	"reverie.jp/reverie/internal/modules/follow/usecase"
	useradapter "reverie.jp/reverie/internal/modules/user/adapter"
	"reverie.jp/reverie/internal/platform/resourcename"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromFollowUserRequest(ctx context.Context, req *connect.Request[followv1.FollowUserRequest]) (usecase.FollowUserInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.FollowUserInput{}, xerrors.ErrUnauthenticated
	}
	customID, err := resourcename.ParseUser(req.Msg.Name)
	if err != nil {
		return usecase.FollowUserInput{}, err
	}
	return usecase.FollowUserInput{
		RequesterID:    requesterID,
		TargetCustomID: customID,
	}, nil
}

func ToFollowUserResponse(output *usecase.FollowUserOutput) *connect.Response[followv1.FollowUserResponse] {
	return connect.NewResponse(&followv1.FollowUserResponse{
		User: useradapter.ToUser(output.View),
	})
}
