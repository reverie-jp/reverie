package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/ulid"
)

func FromListPostRepostsRequest(ctx context.Context, req *connect.Request[postv1.ListPostRepostsRequest]) (usecase.ListPostRepostsInput, ulid.ULID, error) {
	userID, _ := interceptor.UserIDFromContext(ctx)
	return usecase.ListPostRepostsInput{
		AuthorCustomID: req.Msg.AuthorId,
		ShortID:        req.Msg.ShortId,
		PageToken:      req.Msg.PageToken,
		PageSize:       req.Msg.PageSize,
	}, userID, nil
}

func ToListPostRepostsResponse(outputs []*usecase.PostOutput) *connect.Response[postv1.ListPostRepostsResponse] {
	return connect.NewResponse(&postv1.ListPostRepostsResponse{
		Posts:         ToPosts(outputs),
		NextPageToken: NextPageToken(outputs),
	})
}
