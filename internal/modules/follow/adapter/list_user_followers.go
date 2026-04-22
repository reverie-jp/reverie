package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	followv1 "reverie.jp/reverie/internal/gen/pb/follow/v1"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/modules/follow/usecase"
	useradapter "reverie.jp/reverie/internal/modules/user/adapter"
	"reverie.jp/reverie/internal/platform/resourcename"
)

func FromListUserFollowersRequest(ctx context.Context, req *connect.Request[followv1.ListUserFollowersRequest]) (usecase.ListUserFollowersInput, error) {
	customID, err := resourcename.ParseUser(req.Msg.Name)
	if err != nil {
		return usecase.ListUserFollowersInput{}, err
	}
	input := usecase.ListUserFollowersInput{
		TargetCustomID: customID,
		PageSize:       req.Msg.PageSize,
		PageToken:      req.Msg.PageToken,
	}
	if requesterID, ok := interceptor.UserIDFromContext(ctx); ok {
		input.RequesterID = requesterID
	}
	return input, nil
}

func ToListUserFollowersResponse(output *usecase.ListUserFollowersOutput) *connect.Response[followv1.ListUserFollowersResponse] {
	users := make([]*userv1.User, 0, len(output.Views))
	for _, v := range output.Views {
		if u := useradapter.ToUser(v); u != nil {
			users = append(users, u)
		}
	}
	return connect.NewResponse(&followv1.ListUserFollowersResponse{
		Users:         users,
		NextPageToken: output.NextPageToken,
	})
}
