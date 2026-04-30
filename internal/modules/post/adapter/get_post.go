package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/ulid"
)

func FromGetPostRequest(ctx context.Context, req *connect.Request[postv1.GetPostRequest]) (usecase.GetPostInput, ulid.ULID, error) {
	userID, _ := interceptor.UserIDFromContext(ctx)
	return usecase.GetPostInput{
		AuthorCustomID: req.Msg.AuthorId,
		ShortID:        req.Msg.ShortId,
	}, userID, nil
}

func ToGetPostResponse(output *usecase.PostOutput) *connect.Response[postv1.GetPostResponse] {
	return connect.NewResponse(&postv1.GetPostResponse{Post: ToPost(output)})
}
