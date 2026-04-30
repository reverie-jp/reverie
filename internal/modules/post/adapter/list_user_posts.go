package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/ulid"
)

func FromListUserPostsRequest(ctx context.Context, req *connect.Request[postv1.ListUserPostsRequest]) (usecase.ListUserPostsInput, ulid.ULID, error) {
	userID, _ := interceptor.UserIDFromContext(ctx)
	return usecase.ListUserPostsInput{
		UserID:    req.Msg.UserId,
		PageToken: req.Msg.PageToken,
		PageSize:  req.Msg.PageSize,
	}, userID, nil
}

func ToListUserPostsResponse(outputs []*usecase.PostOutput) *connect.Response[postv1.ListUserPostsResponse] {
	return connect.NewResponse(&postv1.ListUserPostsResponse{
		Posts:         ToPosts(outputs),
		NextPageToken: NextPageToken(outputs),
	})
}
