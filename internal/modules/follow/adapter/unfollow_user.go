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

func FromUnfollowUserRequest(ctx context.Context, req *connect.Request[followv1.UnfollowUserRequest]) (usecase.UnfollowUserInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.UnfollowUserInput{}, xerrors.ErrUnauthenticated
	}
	customID, err := resourcename.ParseUser(req.Msg.Name)
	if err != nil {
		return usecase.UnfollowUserInput{}, err
	}
	return usecase.UnfollowUserInput{
		RequesterID:    requesterID,
		TargetCustomID: customID,
	}, nil
}

func ToUnfollowUserResponse(output *usecase.UnfollowUserOutput) *connect.Response[followv1.UnfollowUserResponse] {
	return connect.NewResponse(&followv1.UnfollowUserResponse{
		User: useradapter.ToUser(output.View),
	})
}
