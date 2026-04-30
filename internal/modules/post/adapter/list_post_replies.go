package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/ulid"
)

func FromListPostRepliesRequest(ctx context.Context, req *connect.Request[postv1.ListPostRepliesRequest]) (usecase.ListPostRepliesInput, ulid.ULID, error) {
	userID, _ := interceptor.UserIDFromContext(ctx)
	return usecase.ListPostRepliesInput{
		AuthorCustomID: req.Msg.AuthorId,
		ShortID:        req.Msg.ShortId,
		PageToken:      req.Msg.PageToken,
		PageSize:       req.Msg.PageSize,
	}, userID, nil
}

func ToListPostRepliesResponse(outputs []*usecase.PostOutput) *connect.Response[postv1.ListPostRepliesResponse] {
	return connect.NewResponse(&postv1.ListPostRepliesResponse{
		Posts:         ToPosts(outputs),
		NextPageToken: NextPageToken(outputs),
	})
}
