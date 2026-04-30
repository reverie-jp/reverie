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

func FromCreatePostRequest(ctx context.Context, req *connect.Request[postv1.CreatePostRequest]) (usecase.CreatePostInput, ulid.ULID, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.CreatePostInput{}, ulid.ULID{}, xerrors.ErrUnauthenticated
	}
	return usecase.CreatePostInput{
		Text:      req.Msg.Text,
		ReplyToPostID: req.Msg.ReplyToPostId,
		RepostPostID:  req.Msg.RepostPostId,
	}, userID, nil
}

func ToCreatePostResponse(output *usecase.PostOutput) *connect.Response[postv1.CreatePostResponse] {
	return connect.NewResponse(&postv1.CreatePostResponse{Post: ToPost(output)})
}
