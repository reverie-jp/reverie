package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromListUserLikedPostsRequest(ctx context.Context, req *connect.Request[postv1.ListUserLikedPostsRequest]) (usecase.ListUserLikedPostsInput, ulid.ULID, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.ListUserLikedPostsInput{}, ulid.ULID{}, xerrors.ErrUnauthenticated
	}
	return usecase.ListUserLikedPostsInput{
		UserID:    req.Msg.UserId,
		PageToken: req.Msg.PageToken,
		PageSize:  req.Msg.PageSize,
	}, userID, nil
}

func ToListUserLikedPostsResponse(outputs []*usecase.PostOutput) *connect.Response[postv1.ListUserLikedPostsResponse] {
	return connect.NewResponse(&postv1.ListUserLikedPostsResponse{
		Posts:         ToPosts(outputs),
		NextPageToken: NextPageToken(outputs),
	})
}
